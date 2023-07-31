# HF WSS Convo

## Description
A WSS integration to communicate with the convo from your terminal using the browser as a proxy.

## Usage
1. Clone the repo

2. Create a self-signed certificate for upgrading WS to WSS
```bash
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes
```

3. Add the tampermonkey script to your browser
 * *Remember to set the correct username at the top of the script*

4. Run the server
```bash
go run main.go
```

5. Open the browser and go to the convo

# TODO

- [ ] Colorize the usernames depending on groups and roles
- [ ] Show logged in users
- [ ] Ask for username on startup
- [ ] Refresh if connection to the convo is lost

- [x] Add a way to send messages to the convo from the terminal
