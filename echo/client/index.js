"use strict";

{
  const ws = new WebSocket("ws://localhost:8080/ws");
  ws.onopen = () => {
    console.log("Websocket connection established.");
  };

  const messageList = document.getElementById("messages");
  ws.onmessage = (event) => {
    const message = event.data;
    messageList.innerHTML += `<li>Echo: ${message}</li>`;
  };
  ws.onerror = (event) => {
    console.error("Websocket error: ", event);
  };
  ws.onclose = () => {
    console.log("Websocket connection closed.");
  };

  document.getElementById("sendBtn").addEventListener("click", (ev) => {
    ev.target.innerHTML = "Sending...";
    ev.target.disabled = true;

    const inputElm = document.getElementById("input");
    const message = inputElm.value;
    ws.send(message);
    messageList.innerHTML += `<li>You: ${message}</li>`;

    ev.target.innerHTML = "Send";
    ev.target.disabled = false;
  });
}
