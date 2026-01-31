package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"

	// Import your API definitions
	cortexv1alpha1 "github.com/KGEmmanuel/KubeAegis/api/v1alpha1"
)

var (
	scheme    = runtime.NewScheme()
	k8sClient client.Client
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(cortexv1alpha1.AddToScheme(scheme))
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "aegis",
		Short: "Aegis AI Operator CLI",
	}

	var planCmd = &cobra.Command{
		Use:   "plan [intent]",
		Short: "Ask the AI to generate a plan",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			intent := args[0]
			fmt.Printf("🛡️  Aegis is thinking about: %q...\n", intent)
			createPlan(intent)
		},
	}

	rootCmd.AddCommand(planCmd)
	initClient()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initClient() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	kubeconfig := filepath.Join(home, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		fmt.Printf("❌ Could not find Kubeconfig: %v\n", err)
		os.Exit(1)
	}

	k8sClient, err = client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		fmt.Printf("❌ Could not create K8s client: %v\n", err)
		os.Exit(1)
	}
}

func createPlan(intent string) {
	ctx := context.Background()

	// 1. Create the Object
	req := &cortexv1alpha1.PlanRequest{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "cli-req-",
			Namespace:    "default",
		},
		Spec: cortexv1alpha1.PlanRequestSpec{
			Intent:  intent,
			Planner: "autogen",
		},
	}

	if err := k8sClient.Create(ctx, req); err != nil {
		fmt.Printf("❌ Failed to send request: %v\n", err)
		return
	}

	fmt.Printf("✅ Request Sent! ID: %s\n", req.Name)
	fmt.Print("⏳ Waiting for the Brain to generate code...")

	// 2. Poll for results (Wait up to 15 seconds)
	for i := 0; i < 15; i++ {
		var latest cortexv1alpha1.PlanRequest
		err := k8sClient.Get(ctx, types.NamespacedName{
			Name:      req.Name,
			Namespace: req.Namespace,
		}, &latest)

		if err != nil {
			fmt.Printf("\n❌ Error fetching status: %v\n", err)
			return
		}

		// 3. Check Success State
		if latest.Status.Phase == "Validating" || latest.Status.Phase == "Completed" {
			fmt.Println("\n✨ Plan Generated Successfully:\n")
			fmt.Println("---")
			
			// Clean up Markdown formatting
			cleanYaml := strings.ReplaceAll(latest.Status.GeneratedPlan, "```yaml", "")
			cleanYaml = strings.ReplaceAll(cleanYaml, "```", "")
			fmt.Println(cleanYaml)
			fmt.Println("---")

			// 4. Interactive Menu (Now with Refine!)
			fmt.Print("\n❓ What would you like to do? [A]pply / [R]efine / [Q]uit: ")
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				choice := strings.ToUpper(strings.TrimSpace(scanner.Text()))
				
				switch choice {
				case "A":
					applyYaml(cleanYaml)
				case "R":
					// --- REFINEMENT LOGIC ---
					fmt.Print("\n📝 Enter your feedback (what should change?): ")
					if scanner.Scan() {
						feedback := scanner.Text()
						fmt.Println("\n🔄 Refining plan based on feedback...")
						
						// Construct a new "Super Prompt" that includes context
						newIntent := fmt.Sprintf(
							"ORIGINAL REQUEST: %s\n\nPREVIOUS PLAN GENERATED:\n%s\n\nUSER FEEDBACK: %s\n\nINSTRUCTION: Regenerate the YAML to address the user feedback.", 
							intent, cleanYaml, feedback,
						)
						
						// Recursive call with the new "Super Prompt"
						createPlan(newIntent) 
					}
				default:
					fmt.Println("👋 Plan discarded.")
				}
			}
			return
		}

		if latest.Status.Phase == "Failed" {
			fmt.Printf("\n❌ AI Failed: %s\n", latest.Status.Reason)
			return
		}

		time.Sleep(1 * time.Second)
		fmt.Print(".")
	}

	fmt.Println("\n⚠️  Timeout waiting for AI. Run 'kubectl get planrequests' to check later.")
}

// applyYaml pipes the generated YAML directly to kubectl apply
func applyYaml(yamlContent string) {
	cmd := exec.Command("kubectl", "apply", "-f", "-")
	cmd.Stdin = strings.NewReader(yamlContent)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("\n🚀 Applying plan to cluster...")
	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ Failed to apply: %v\n", err)
	} else {
		fmt.Println("✅ Deployment Successful!")
	}
}
