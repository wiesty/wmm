# 🎮 Wiesty's Mod Installer 🛠️ <img src="https://img.shields.io/badge/Maintained%3F-no-red.svg"/>

![Mod Installer](https://i.imgur.com/nvEUQon.png)

## 📝 Overview

**Wiesty's Mod Installer** is a terminal-based tool designed to simplify the installation and restoration of mods for Minecraft. It allows you to backup, restore, and install mods with just a few keystrokes, making mod management effortless. 🎉

With this tool, you can:
- 🗂️ Automatically create backups of your `.minecraft` folder.
- 📥 Download and install mods and installers with ease.
- 🔄 Restore your Minecraft folder from any of your previous backups.

## 🚀 Features

- **Automatic Backup Creation**: Before installing any mods, the tool automatically backs up your existing `.minecraft` folder with a timestamp.
- **Backup Management**: Easily restore or delete backups from a list.
- **Mod Installation**: Download and install mods from URLs defined in the `wmm.json` config file.
- **Installer Support**: Run any external mod installer (like Fabric) with no hassle.
- **User-friendly CLI**: Intuitive CLI with arrow key navigation and easy-to-understand prompts.

## 🛠️ How It Works

1. **Select the Path**: The installer automatically detects your default `.minecraft` path, or you can choose a custom one.
2. **Automatic Backup**: Once the path is selected, the tool creates a backup of your `.minecraft` folder. The backup is named with a timestamp (e.g., `.minecraftBACKUP_2024-10-18_15-30-05`).
3. **Download and Run Installer**: The installer defined in the `wmm.json` (for example Fabric) will be downloaded and run. You can choose any external installer to be executed.
4. **Download and Install Mods**: After the installer is finished, the mods defined in the config will be downloaded and unzipped into the `.minecraft` folder.
5. **Restore and Delete Backups**: You can list, restore, or delete backups easily through the menu.

## 🔧 Configuration (`wmm.json`)

Here's a sample `wmm.json` file that defines the hints, installer URLs, and mods:

```json
{
  "hints": [
    "Please select the correct Minecraft version in the installer.",
    "Do not close the terminal while the installer is running."
  ],
  "installurl": [
    "https://your-url-to-installer/some-installer-1.0.1.exe"
  ],
  "modsurl": "https://your-url-to-mods/mods.zip"
}
```


### Explanation:

-   **`hints`**: These are helpful messages that will be displayed to the user during the installation process. You can provide as many hints as you like to guide the user.
-   **`installurl`**: One or more direct URLs pointing to the installer executables. If there are multiple installers, each one will be downloaded and executed one at a time. The URLs must be **direct download links** (e.g., not mediafire). Examples include URLs from file hosting services that allow direct file downloads.
-   **`modsurl`**: A direct URL to the mods package (`.zip` file) that will be downloaded and installed. The `.zip` file should have the exact folder structure of a `.minecraft` folder for the mods to be correctly placed (i.e., it should contain folders like `mods/`, `resourcepacks/`, etc.).

### Creating the `mods.zip` File:

-   The `mods.zip` file should be structured to match the folder layout of the `.minecraft` directory.
    -   Example contents of `mods.zip`:
        ```
        mods/
        ├── ExampleMod1.jar
        ├── ExampleMod2.jar
        resourcepacks/
        ├── ExampleResourcePack.zip
        options.txt
        ``` 
        
-   Make sure the `mods/` folder contains your mods and any other folders (like `resourcepacks/`, `config/`) as needed.

### Important Notes:

-   The **links in the config** must be direct links to the files. Mediafire or any other file-sharing services that don't provide direct links will not work.
-   You can provide **multiple installers** in the `installurl` array if needed. The tool will download and run each installer one by one.

## 🛠️ Installation & Setup

1.  By downloading you agree to the [Disclaimer](https://raw.githubusercontent.com/wiesty/wmm/refs/heads/main/Disclaimer.txt). **Download the `wmm.exe`** from the [Releases](https://github.com/wiesty/wmm/releases/tag/v1).
2.  **Create your `wmm.json` file** with the structure shown above.
3.  **Run the Installer**:
    -   Open a terminal or simply double-click `wmm.exe`.
    -   Follow the prompts to select the path, backup your current `.minecraft` folder, download the installer, and install the mods.

### Example Usage Flow:

1.  **Path Selection**: The installer automatically detects your `.minecraft` path or lets you specify a custom path.
2.  **Backup Creation**: A backup of the `.minecraft` folder is created with a timestamp (e.g., `.minecraftBACKUP_2024-10-18_15-30-05`).
3.  **Installer Execution**: The external installer (like Fabric) is downloaded and run.
4.  **Mods Installation**: The `mods.zip` file is downloaded, and its contents are extracted into your `.minecraft` folder.
5.  **Post-Task Menu**: After completing the installation, you can choose to return to the main menu or exit.

## 💡 Tips

-   Ensure your `wmm.json` config is properly set up before running the installer.
-   If no installer is needed, leave the `installurl` field empty in the `wmm.json`.
-   Backups are created with timestamps, so you'll never lose track of previous states of your `.minecraft` folder.
-   You can easily switch between different mod packs by restoring backups or creating new mod configurations.

**Disclaimer:**

- 💻 See the Disclaimer [here](https://raw.githubusercontent.com/wiesty/wmm/refs/heads/main/Disclaimer.txt).

