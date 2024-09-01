# DisVault

DisVault is a lightweight file management solution that leverages Discord servers to store, organize, and manage files. **Caution: This project is in alpha stage**, and it's not recommended to upload large volumes of files to your Discord server as it may exceed Discord's limitations.

## 🚧 **Project Status: Alpha**

**⚠️ Warning:** DisVault is currently in alpha, meaning it is under active development and may contain bugs or missing features. Use with caution, especially when uploading sensitive or large amounts of data. 

## 📜 **Overview**

DisVault allows you to:

- **Upload:** Automatically Split files into manageable chunks and upload them to your Discord server.
- **Download:** Retrieve your files seamlessly by assembling the chunks back together.
- **List:** View your files in an organized manner, grouped for easier access and management.
- **Delete:** Remove files from the server when they're no longer needed.
- **Groups:** Assign files to groups for easier categorization and searching.

### ⚡ Quick Start

1. **Clone the Repository**

   ```bash
   git clone https://github.com/yourusername/disvault.git
   cd disvault
   ```

2. **Set Up Your Environment**

   - Ensure you have Go installed (version 1.20 or later).
   - Set up your Discord bot with the appropriate permissions.

3. **Build and Run**

   ```bash
   go build -o disvault
   ./disvault --help
   ```

4. **Upload Files**

   Use the upload command to add files:

   ```bash
   ./disvault upload --file yourfile.txt
   ```

## 📋 **Features**

- **Lightweight**: Minimal setup and dependencies.
- **Discord Integration**: Uses Discord channels for file storage.
- **Searchable**: Easily search and filter files using various flags.
- **Categorization**: Group files to keep everything organized.

## 🚀 **Usage**

```bash
disvault [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  delete      Delete files using their IDs
  download    Download files using their IDs
  help        Help about any command
  list        List the uploaded files
  upload      Upload a file by splitting it into chunks and registering it in the database
  version     Print the version number of DisVault
```

## ⚠️ **Caution**

- **Discord Limitations**: Uploading a large number of files or very large files can exceed Discord’s storage limitations and could get your bot rate-limited or banned.
- **Data Integrity**: This is an alpha release, and while we strive for data integrity, please keep backups of important files elsewhere.

## 💡 **Contributing**

We welcome contributions! Please fork the repository and submit pull requests for new features or bug fixes. For major changes, please open an issue first to discuss what you would like to change.

## 📝 **TODO**

- [ ] Improve error handling and logging.
- [ ] Add support for more file types and formats.
- [ ] Enhance file search functionality with more filters.
- [ ] Develop a web-based interface for easier file management.
- [ ] Implement user authentication and permissions.

## 🛠️ **Built With**

- [Go](https://golang.org/) - Programming Language
- [SQLite](https://sqlite.org/) - Database
- [Discord API](https://discord.com/developers/docs/intro) - For file storage and retrieval

## 📄 **License**

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## ❤️ **Support**

If you find DisVault useful, please give it a ⭐ on GitHub and share it with your friends and colleagues!