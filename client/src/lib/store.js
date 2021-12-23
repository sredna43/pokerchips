import { writable } from 'svelte/store'
import { browser } from '$app/env'

export const messageStore = writable('');
export const serverHost = '192.168.0.12:8081'

let socket;

if (browser) {
    socket = new WebSocket('ws://' + serverHost + '/ws');

    socket.addEventListener('open', (event) => {
        console.log('opened websocket');
    });

    socket.addEventListener('message', (event) => {
        messageStore.set(event.data);
    });
}

export const sendMessage = (message) => {
    if (socket.readyState <= 1 && browser) {
        console.log('sending message');
        socket.send(message);
    }
}