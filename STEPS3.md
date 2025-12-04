Super ! 🎉 Content de savoir que tout fonctionne parfaitement. Voici la version finale du **README** pour ton projet `archetype-cli`, reflétant le flux complet et le support des flags `--install`, `--test` et `--custom`.

---

# Archetype CLI

`archetype-cli` est une application en **Go** pour installer et générer facilement votre archetype Maven localement, sans passer par un dépôt central comme Nexus.

## Fonctionnalités

* Vérifie si **Java** et **Maven** sont installés.
* Installe le JAR de l’archetype localement.
* Permet de générer un projet Maven à partir de l’archetype.
* Supporte des valeurs personnalisées pour :

  * `groupId`
  * `artifactId`
  * `version`
  * `package`
* Supporte les flags CLI :

  * `--install` : installe le JAR de l’archetype
  * `--test` : génère automatiquement un projet Maven pour tester
  * `--custom` : personnalise les valeurs Maven via un prompt interactif

## Installation

1. Clonez le projet et build l’application Go :

```bash
git clone <repo-url>
cd archetype-cli
go build -o archetype-cli main.go
```

2. Assurez-vous que votre JAR de l’archetype est présent :

```
factory/usine.jar
```

## Utilisation

### 1️⃣ Vérification des outils requis

```bash
./archetype-cli
```

Affiche les versions installées de Java et Maven.

---

### 2️⃣ Installer le JAR localement

```bash
./archetype-cli --install
```

Installe le JAR dans le repository Maven local (`~/.m2/repository`).

---

### 3️⃣ Tester la génération automatiquement avec les valeurs par défaut

```bash
./archetype-cli --install --test
```

Génère un projet Maven en utilisant l’archetype installé avec les valeurs par défaut :

```text
GroupId: baobao
ArtifactId: pi-zb
Version: 1.0-SNAPSHOT
Package: sn.cbao
```

---

### 4️⃣ Personnaliser les valeurs Maven

```bash
./archetype-cli --install --custom
```

* L’utilisateur est invité à renseigner :

  * `GroupId` (par défaut : baobao)
  * `ArtifactId` (par défaut : pi-zb)
  * `Version` (par défaut : 1.0-SNAPSHOT)
  * `Package` (par défaut : sn.cbao)
* L’utilisateur peut choisir de générer le projet immédiatement ou copier la commande Maven proposée pour exécution ultérieure.

Exemple de commande générée :

```bash
mvn archetype:generate \
  -DarchetypeCatalog=local \
  -DarchetypeGroupId=com.votreorganisation.archetypes \
  -DarchetypeArtifactId=starter-kit-archetype \
  -DarchetypeVersion=0.0.1-SNAPSHOT \
  -DgroupId=titi \
  -DartifactId=poutpour \
  -Dversion=1.0-SNAPSHOT \
  -Dpackage=kd.wassadou \
  -DinteractiveMode=false
```

---

### 5️⃣ Flags disponibles

| Flag        | Description                                         |
| ----------- | --------------------------------------------------- |
| `--install` | Installe le JAR de l’archetype dans le repo local   |
| `--test`    | Teste la génération d’un projet Maven immédiatement |
| `--custom`  | Permet de renseigner des valeurs personnalisées     |

---

### 6️⃣ Structure du projet

```
archetype-cli/
│
├── factory/
│   └── usine.jar
├── main.go
└── README.md
```

---

### 7️⃣ Notes

* Compatible avec Go 1.21+.
* Java 17 et Maven 3.9+ recommandés.
* Les valeurs par défaut sont utilisées si l’utilisateur laisse les champs vides.

---

Si tu veux, je peux aussi te **fournir le fichier `.gitignore` optimisé pour ce projet Go + Maven** pour que ton repo soit propre et prêt à partager.

Veux‑tu que je fasse ça ?
