import { writable } from 'svelte/store'
import { browser } from '$app/env'

export const messageStore = writable('');
export const socketState = writable('');
export const serverHost = 'localhost:8081'

let socket;

if (browser) {
    socket = new WebSocket('ws://' + serverHost + '/ws');

    socket.addEventListener('open', (event) => {
        console.log('opened websocket');
        socketState.set('connected');
    });

    socket.addEventListener('message', (event) => {
        messageStore.set(event.data);
    });
}

export const sendMessage = (message) => {
    if (socket.readyState === WebSocket.OPEN && browser) {
        socket.send(message);
    }
}