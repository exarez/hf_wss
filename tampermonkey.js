// ==UserScript==
// @name     HackForums Chat WSS
// @author   You
// @description Communicate with the convo using the browser as proxy
// @version  1
// @grant    none
// @match    https://hackforums.net/convo.php
// ==/UserScript==

const username = "Ezocain";

(function () {
	"use strict";

	// We need the Convo object, so we'll wait for it to appear
	function checkAuthentication() {
		if (Convo.group_users !== undefined) {
			clearInterval(checkInterval);
			startScript();
		} else {
			console.log("Waiting for authentication...");
		}
	}

	// Start checking for the element every second
	var checkInterval = setInterval(checkAuthentication, 1000);

	// Define a function that starts your script
	function startScript() {
		console.log("Starting script...");

		// Buttons etc.
		const convoControls = document.getElementById("convoControlsRow");
		const sendBtn = convoControls.querySelector("input[name='submit_button']");
		const messageBox = convoControls.querySelector("textarea[name='comment']");

		const testMessage = () => {
			setInput("testing");
		};

		const setInput = (input) => {
			if (input === "") return;

			messageBox.value = input;
			send();
		};

		const send = () => {
			if (messageBox.value === "") {
				console.log("No input to send - returning");
				return;
			}
			// window.Convo.sendMessage();
			sendBtn.click();
		};

		console.log("Input found: ", messageBox);

		console.log("Initiating websockets");
		const WEB_SOCKET_URL = "ws://127.0.0.1:3333"; // replace with your WebSocket server's URL
		const socket = new WebSocket(WEB_SOCKET_URL);
		let lastUid = null;
		let lastUsername = null;

		// Connection opened
		socket.addEventListener("open", (event) => {
			console.log("Connection opened");
			if (event.data) {
				console.log(event.data);
			}
			// socket.send('Hello Server!');
		});

		socket.onmessage = (event) => {
			console.log("Message from server ", event.data);
			let data = JSON.parse(event.data);
			if (data.username == username && data.content) {
				setInput(data.content);
			} else {
				console.log(
					"Error sending content. Content blank, or make sure you set the correct username in the tampermonkey script."
				);
			}
		};

		// Connection closed
		socket.addEventListener("close", (event) => {
			console.log("Connection closed");
		});

		// Connection error
		socket.addEventListener("error", (event) => {
			console.log("WebSocket error: ", event);
		});

		function fetchLatestMessages() {
			let allMessages = Array.from(
				document.querySelectorAll(".message-convo-left, .message-convo-right")
			);
			let newMessages = [];
			for (let i = allMessages.length - 1; i >= 0; i--) {
				if (allMessages[i].getAttribute("data-uid") !== lastUid) {
					newMessages.unshift(allMessages[i]);
				} else {
					break;
				}
			}
			return newMessages;
		}

		// let convo = {...Convo};

		// console.log("Users: ", convo);
		function matchUserToGroup(username) {
			// Iterate over each user in the users object
			for (let uid in Convo.group_users) {
				// Check if this user's username matches the one we're looking for
				if (Convo.group_users[uid].username === username) {
					// If it does, return this user's group
					return Convo.group_users[uid].usergroup;
				}
			}

			return null;
		}

		function sendViaWebSocket(username, message, usergroup) {
			console.log(
				`Sending user with username ${username} in usergroup ${usergroup} to the server...`
			);
			if (socket.readyState === WebSocket.OPEN) {
				socket.send(
					JSON.stringify({
						username: username,
						content: message,
						usergroup: usergroup,
					})
				);
			} else {
				console.error("WebSocket is not open");
			}
		}

		function sendUserList(users) {
			if (socket.readyState === WebSocket.OPEN) {
				socket.send(
					JSON.stringify({
						users: users,
					})
				);
			}
		}
		// testMessage();

		setInterval(() => {
			let latestMessages = fetchLatestMessages();
			for (let latestMessage of latestMessages) {
				let currentUid = latestMessage.getAttribute("data-uid");

				let usernameDiv = latestMessage.querySelector(".message-bubble-left a");
				let usernameFromAvatar = latestMessage.querySelector(
					".mirum-card-profile-info-username a span"
				);

				let username;
				if (latestMessage.classList.contains("message-convo-right")) {
					username = "Ezocain";
				} else if (usernameDiv) {
					username = usernameDiv.textContent.trim();
					lastUsername = username; // Save the username for further use
				} else if (usernameFromAvatar) {
					username = usernameFromAvatar.textContent.trim();
					lastUsername = username; // Save the username for further use
				} else if (currentUid === lastUid && lastUsername) {
					username = lastUsername;
				} else {
					console.log("Could not find username for message: ", latestMessage);
					continue;
				}

				let messageDiv = latestMessage.querySelector(".message-bubble-message");
				if (messageDiv) {
					let message = messageDiv.textContent.trim();
					let usergroup = matchUserToGroup(username);

					console.log(
						`New message from ${username}\(${usergroup}\): ${message}`
					);
					sendViaWebSocket(username, message, usergroup);
				}

				lastUid = currentUid;
			}
		}, 100);

		// For a later feature for showing online members
		// setInterval(() => {
		//     console.log("Sending users to client...")
		//     let users = {...Convo.group_users}
		//     if (users.length > 0) {
		//         sendUserList(
		//             users
		//         )}
		// }, 10000);
	}
})();
