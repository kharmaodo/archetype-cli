Parfait, tu veux maintenant ajouter une **option `--custom`** pour que l’utilisateur puisse renseigner les valeurs dynamiques (`groupId`, `artifactId`, `version`, `package`) de manière interactive, avec des valeurs par défaut si l’utilisateur appuie simplement sur Entrée.

Voici comment réécrire `main.go` pour gérer ce cas :

```go
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// Flags CLI
	install := flag.Bool("install", false, "Installer le jar localement")
	test := flag.Bool("test", false, "Tester la génération d'un projet Maven depuis l'archetype")
	custom := flag.Bool("custom", false, "Permet de renseigner les valeurs groupId, artifactId, version, package")
	groupId := flag.String("groupId", "baobao", "GroupId du projet généré")
	artifactId := flag.String("artifactId", "pi-zb", "ArtifactId du projet généré")
	version := flag.String("version", "1.0-SNAPSHOT", "Version du projet généré")
	packageName := flag.String("package", "sn.cbao", "Package du projet généré")
	flag.Parse()

	fmt.Println("🔍 Vérification des outils requis...")

	javaVersion, javaOk := checkCommand("java", "-version")
	mvnVersion, mvnOk := checkCommand("mvn", "-v")

	if !javaOk || !mvnOk {
		fmt.Println("❌ Java ou Maven non installé. Veuillez installer Java et Maven.")
		return
	}

	fmt.Println("✅ Java:", javaVersion)
	fmt.Println("✅ Maven:", mvnVersion)

	jarPath := "factory/starter-kit-archetype-0.0.1-SNAPSHOT.jar"
	if _, err := os.Stat(jarPath); os.IsNotExist(err) {
		fmt.Printf("❌ Fichier JAR introuvable : %s\n", jarPath)
		return
	}

	if *install {
		fmt.Println("⚙️ Installation du jar en local...")
		err := runCommand("mvn", "install:install-file",
			"-Dfile="+jarPath,
			"-DgroupId=com.votreorganisation.archetypes",
			"-DartifactId=starter-kit-archetype",
			"-Dversion=0.0.1-SNAPSHOT",
			"-Dpackaging=jar",
		)
		if err != nil {
			fmt.Println("❌ Échec de l'installation :", err)
			return
		}
		fmt.Println("✅ Installation réussie !")
		fmt.Println("\nVous pouvez maintenant exécuter : mvn archetype:generate -DarchetypeCatalog=local")
	}

	reader := bufio.NewReader(os.Stdin)

	if *custom {
		fmt.Println("\n⚡ Personnalisation des valeurs Maven")

		fmt.Printf("GroupId (%s) : ", *groupId)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input != "" {
			*groupId = input
		}

		fmt.Printf("ArtifactId (%s) : ", *artifactId)
		input, _ = reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input != "" {
			*artifactId = input
		}

		fmt.Printf("Version (%s) : ", *version)
		input, _ = reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input != "" {
			*version = input
		}

		fmt.Printf("Package (%s) : ", *packageName)
		input, _ = reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input != "" {
			*packageName = input
		}
	}

	if *test {
		fmt.Print("\n💡 Voulez-vous générer un projet Maven avec ces paramètres ? (y/n) : ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)
		if strings.ToLower(choice) == "y" {
			fmt.Println("\n⚙️ Génération du projet Maven...")
			err := runCommand("mvn", "archetype:generate",
				"-DarchetypeCatalog=local",
				"-DarchetypeGroupId=com.votreorganisation.archetypes",
				"-DarchetypeArtifactId=starter-kit-archetype",
				"-DarchetypeVersion=0.0.1-SNAPSHOT",
				"-DgroupId="+*groupId,
				"-DartifactId="+*artifactId,
				"-Dversion="+*version,
				"-Dpackage="+*packageName,
				"-DinteractiveMode=false",
			)
			if err != nil {
				fmt.Println("❌ Échec de la génération :", err)
				return
			}
			fmt.Println("✅ Projet Maven généré avec succès !")
		} else {
			fmt.Println("\nℹ️ Pour générer manuellement, utilisez la commande suivante :")
			fmt.Printf(`mvn archetype:generate \
  -DarchetypeCatalog=local \
  -DarchetypeGroupId=com.votreorganisation.archetypes \
  -DarchetypeArtifactId=starter-kit-archetype \
  -DarchetypeVersion=0.0.1-SNAPSHOT \
  -DgroupId=%s \
  -DartifactId=%s \
  -Dversion=%s \
  -Dpackage=%s \
  -DinteractiveMode=false
`, *groupId, *artifactId, *version, *packageName)
		}
	}
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func checkCommand(name string, arg string) (string, bool) {
	cmd := exec.Command(name, arg)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", false
	}
	version := strings.TrimSpace(string(out))
	version = strings.Split(version, "\n")[0]
	return version, true
}
```

---

✅ **Fonctionnalités ajoutées avec `--custom`** :

1. L’utilisateur peut renseigner interactif les valeurs de :

   * `groupId`
   * `artifactId`
   * `version`
   * `package`
     avec des valeurs par défaut si rien n’est saisi.

2. Compatible avec `--install` et `--test` simultanément.

3. Maintient la génération automatisée pour scripts CI/CD si aucun `--custom`.

---

Si tu veux, je peux te faire un **diagramme d’utilisation et un guide d’exécution CLI complet**, montrant tous les scénarios possibles (`./archetype-cli`, `--install`, `--test`, `--custom`).

Veux‑tu que je fasse ça ?
