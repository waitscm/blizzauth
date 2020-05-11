# blizzauth
Simple way to get your Blizzard API client request token from your stored client keys in `$HOME/.blizzard`

## Setup
Make a .blizzard directory in your $HOME folder

```
mkdir $HOME/.blizzard
```

Create the following files:

```
echo "<your api client key" > ~/.blizzard/your_api_name.id
echo "<your secret id key" > ~/.blizzard/you_api_name.secret
chmod 400 ~/.blizzard/you_api_name.secret
```

## Usage

```
auth, err := blizzauth.GetAuth("your_api_name")
...
// When you need to use an access token
token, _ := auth.GetAccessToken()
```
