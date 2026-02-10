"use strict";

(function () {
  const messagesBox = document.getElementById("messagesBox");
  const messageInput = document.getElementById("messageInput");
  const sendButton = document.getElementById("sendButton");

  function addMessage(username, text, isOwn) {
    const messageDiv = document.createElement("div");
    messageDiv.className = `message ${isOwn ? "own" : "other"}`;

    messageDiv.innerHTML = `
                <div class="message-content">
                    <div class="message-header">${username}</div>
                    <div class="message-bubble">${text}</div>
                </div>
            `;

    messagesBox.appendChild(messageDiv);
    messagesBox.scrollTop = messagesBox.scrollHeight;
  }

  const ws = new WebSocket("ws://localhost:8080/ws");
  ws.onopen = () => {
    console.log("Websocket connection established");
  };
  ws.onerror = (ev) => {
    console.log("Websocket error: ", ev);
  };
  ws.onmessage = (ev) => {
    const message = ev.data;
    addMessage("Unknown", message, false);
  };

  function sendMessage() {
    const msg = messageInput.value;
    messageInput.value = "";
    ws.send(msg);
    addMessage("You", msg, true);
  }

  sendButton.addEventListener("click", sendMessage);

  messageInput.addEventListener("keypress", (e) => {
    if (e.key === "Enter") {
      sendMessage();
    }
  });
})();
