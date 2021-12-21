import { Writable, writable } from 'svelte/store';

interface Player {
    is_host: boolean;
    spot: number;
    name: string;
    chips: number;
    folded: boolean;
}

interface GameState {
    playing: boolean;
    players: [Player];
    turn: number;
    dealer: number;
    pot: number;
}

interface ErrorMsg {
    cause: string;
    message: string;
}

interface ServerResponse {
    message: string;
    player: Player;
    game_state: GameState;
    error: ErrorMsg;
}

let blankPlayer: Player = {
    is_host: false,
    name: '',
    spot: 0,
    chips: 0,
    folded: false,
}

const messageStore: Writable<ServerResponse> = writable({
    message: '',
    player: blankPlayer,
    game_state: {
        playing: false,
        players: [blankPlayer],
        turn: 0,
        dealer: 0,
        pot: 0,
    },
    error: {
        cause: '',
        message: '',
    }
});
const lobbyId: Writable<string> = writable('');

const socket: WebSocket = new WebSocket('ws://localhost:8081/ws');

socket.addEventListener('open', (event) => {
    console.log("websocket open");
});

socket.addEventListener('message', (event) => {
    messageStore.set(event.data);
    console.log(event.data);
})

const sendMessage = (message) => {
    if (socket.readyState <= 1) {
        socket.send(message);
    }
}

export default {
    subscribe: messageStore.subscribe,
    sendMessage
}