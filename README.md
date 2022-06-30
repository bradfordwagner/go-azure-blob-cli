# Go Azure Blob CLI
- lightweight cli to interface with azure blob storage

## Required Environment Variables
```
export AZURE_STORAGE_ACCOUNT_NAME= # storage account name 
export AZURE_STORAGE_ACCOUNT_KEY=  # storage account access key
```

## Commands
```bash
abc container create
abc container delete
abc container list
abc blob delete
abc blob download # download everything under a path
abc blob list   # -l -c
abc blob upload # -f -c -t - used for single file uploads
```

## Install Completions
```bash
abc completion zsh > /usr/local/share/zsh/site-functions/_abc
```
