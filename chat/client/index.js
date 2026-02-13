"use strict";

(function () {
  const nameScreen = document.getElementById("nameScreen");
  const chatContainer = document.getElementById("chatContainer");
  const joinBtn = document.getElementById("joinButton");
  const nameInput = document.getElementById("nameInput");
  const nameError = document.getElementById("nameError");

  const messagesBox = document.getElementById("messagesBox");
  const messageInput = document.getElementById("messageInput");
  const sendButton = document.getElementById("sendButton");

  localStorage.removeItem("username");

  nameInput.onkeydown = () => {
    if (nameError.classList.contains("hidden")) return;

    nameError.innerText = "";
    nameError.classList.add("hidden");
  };

  joinBtn.onclick = () => {
    const name = nameInput.value;

    if (name.trim() === "") {
      nameError.innerText = "Chat name can't be empty!";
      nameError.classList.remove("hidden");
      return;
    }

    joinChat(name);
  };

  let ws;

  function joinChat(name) {
    ws = new WebSocket(`ws://localhost:8080/ws?name=${name}`);
    ws.onopen = () => {
      const timeout = setTimeout(() => {
        if (ws.readyState === 1) {
          console.log("Websocket connection established");
          localStorage.setItem("username", nameInput.value);
          nameInput.value = "";

          if (!nameScreen.classList.contains("hidden")) {
            nameScreen.classList.add("hidden");
          } else {
            console.error("nameScreen already contains hidden class");
          }

          if (chatContainer.classList.contains("hidden")) {
            chatContainer.classList.remove("hidden");
          } else {
            console.error("chatContainer does not contain hidden class");
          }

          clearTimeout(timeout);
        }
      }, 10);
    };
    ws.onerror = (ev) => {
      console.log("Websocket error: ", ev);
    };
    ws.onmessage = (ev) => {
      const data = ev.data;

      if (!data) {
        console.error("message data is empty.");
        return;
      }

      try {
        const parsedData = JSON.parse(data);
        if (!parsedData.name) {
          throw Error("parsed data does not contain username");
        }
        if (!parsedData.message) {
          throw Error("parsed data does not contain message");
        }

        addMessage(parsedData.name, parsedData.message, false);
      } catch (err) {
        console.error(err);
      }
    };
  }

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

  sendButton.addEventListener("click", sendMessage);
  messageInput.addEventListener("keypress", (e) => {
    if (e.key === "Enter") {
      sendMessage();
    }
  });

  function sendMessage() {
    const msg = messageInput.value;
    if (msg === "") return;

    const username = localStorage.getItem("username");
    if (!username) {
      console.error("current username missing");
      return;
    }

    if (!ws) {
      console.error("websocket connection is empty");
      return;
    }

    ws.send(
      JSON.stringify({
        name: username,
        message: msg,
      }),
    );
    addMessage("You", msg, true);

    messageInput.value = "";
  }
})();
