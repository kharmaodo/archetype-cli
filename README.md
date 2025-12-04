````markdown
# Archetype CLI - Starter Kit Maven

Ce projet fournit une application en ligne de commande (CLI) en **Golang** pour faciliter l'installation et l'utilisation d'un **archetype Maven Spring Boot** localement, sans passer par Nexus ou tout autre dépôt distant.

---

## 🚀 Objectif

- Installer automatiquement un **JAR d'archetype Maven** localement.
- Vérifier que **Java** et **Maven** sont installés sur la machine.
- Fournir une aide et des commandes pour générer un projet Spring Boot à partir de l'archetype.
- Permettre aux développeurs de personnaliser les valeurs dynamiques pour le `groupId`, `artifactId`, `version` et `package`.

---

## 📦 Prérequis

- **Java JDK 17** ou supérieur
- **Apache Maven 3.9.x** ou supérieur
- **Golang** pour exécuter l'application CLI

---

## ⚙️ Installation et utilisation

1. **Cloner le projet CLI** ou télécharger l’archive.
2. Placer le JAR `starter-kit-archetype-0.0.1-SNAPSHOT.jar` dans le répertoire `factory` du projet Golang.

3. Lancer l’application CLI :

```bash
./archetype-cli
````

L’application fera automatiquement :

* Vérification de Java et Maven.
* Si l’un ou l’autre manque, elle s’arrête avec un message d’erreur.
* Sinon, elle propose d’installer le JAR d’archetype localement via Maven :

```bash
mvn install:install-file \
  -Dfile=starter-kit-archetype-0.0.1-SNAPSHOT.jar \
  -DgroupId=com.votreorganisation.archetypes \
  -DartifactId=starter-kit-archetype \
  -Dversion=0.0.1-SNAPSHOT \
  -Dpackaging=jar
```

---

## 🧩 Génération d’un projet à partir de l’archetype

Après installation, l’utilisateur peut :

1. **Voir l’aide** pour la commande Maven :

```bash
mvn archetype:generate -DarchetypeCatalog=local
```

2. **Générer un projet Maven complet** en renseignant les paramètres dynamiques (avec valeurs par défaut) :

* `groupId` : baobao
* `artifactId` : pi-zb
* `version` : 1.0-SNAPSHOT
* `package` : sn.cbao
* `interactiveMode` : false

Exemple complet :

```bash
mvn archetype:generate \
  -DarchetypeCatalog=local \
  -DarchetypeGroupId=com.votreorganisation.archetypes \
  -DarchetypeArtifactId=starter-kit-archetype \
  -DarchetypeVersion=0.0.1-SNAPSHOT \
  -DgroupId=baobao \
  -DartifactId=pi-zb \
  -Dversion=1.0-SNAPSHOT \
  -Dpackage=sn.cbao \
  -DinteractiveMode=false
```

L’application CLI guidera l’utilisateur pour personnaliser ces valeurs si nécessaire.

---

## ✅ Fonctionnalités principales

* Détection automatique de Java et Maven.
* Installation locale de l’archetype Maven.
* Affichage d’aide et test de génération d’un projet Maven.
* Paramètres dynamiques pour `groupId`, `artifactId`, `version` et `package`.
* Compatible avec les environnements de développement les plus utilisés.

---

## 📂 Structure du projet

```
archetype-cli/
├── factory/                     # Contient le JAR de l'archetype
├── cmd/                         # Code source Go pour l'application CLI
├── go.mod                        # Dépendances et version du projet Go
├── README.md                     # Ce fichier
└── main.go                       # Entrée principale de l'application CLI
```

---

## ✨ Contribution

Les contributions sont les bienvenues !

* Forker le projet
* Créer une branche pour vos modifications
* Envoyer un **Pull Request** après test complet

---

## 📄 Licence

Ce projet est sous licence MIT. Voir le fichier `LICENSE` pour plus de détails.

```

---

Si tu veux, je peux aussi te **préparer une version améliorée** qui inclut **des captures d’écran de la CLI en action** et un **exemple de génération de projet Maven** pour que ce README soit prêt à publier sur GitHub et utilisé par d’autres développeurs.  

Veux‑tu que je fasse ça ?
```

