
# Backend for Listing files

You can use this mini-service for listing files of dir on home-server (for example).

For run this app directly run command:
```sh
go run *.go
```

## Endpoints

### {your_server_adress}**/list**
Get list of files in specified directory. 

Parametrs: _dir_ :string - path to folder

Returns: 
{
    Id: string,
    Name: string,
    IsDir: boolean
}

### {your_server_adress}**/delete**
Delete files in current directory. 

Parametrs: 

- _isDir_ :boolean ; 

- _name_ :string - fileName

### {your_server_adress}**/rename**
Rename file. 

Parametrs: 

- _name_ :string - full/path/to/file

- _newName_ :string - new fileName

### {your_server_adress}**/copy**
Copy files. 

Parametrs: 

- _newDir_ :string - full/path/to/new/location

- _dir_ :string - full/path/to/current/location

- _name_ :string - fileName

### {your_server_adress}**/upload**
Upload files to specified directory. 

Parametrs: _dir_ :string - full/path/to/current/dir

### {your_server_adress}**/download**
Download files. 

Parametrs: 

- _name_ :string - full/path/to/file
- _isDir_ :boolean


