// ==UserScript==
// @name     HackForums Chat WSS
// @author   You
// @description Communicate with the convo using the browser as proxy
// @version  1
// @grant    none
// @match    https://hackforums.net/convo.php
// ==/UserScript==

(function () {
	"use strict";

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

	function sendViaWebSocket(username, message) {
		if (socket.readyState === WebSocket.OPEN) {
			socket.send(
				JSON.stringify({
					username: username,
					content: message,
				})
			);
		} else {
			console.error("WebSocket is not open");
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
				console.log(`New message from ${username}: ${message}`);
				sendViaWebSocket(username, message);
			}

			lastUid = currentUid;
		}
	}, 100);
})();
