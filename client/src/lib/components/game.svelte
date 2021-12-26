<script>
    import Button from './button.svelte';
    import Input from './input.svelte';
    import Error from './error.svelte';
    import Players from './players.svelte';
    import Settings from './settings.svelte';
    import { onMount } from 'svelte';
    import { serverHost, sendMessage, messageStore, socketState } from '../store.js';

    var lobbyId;
    var myName;
    var state;
    var errorMessage;
    var allPlayers;
    var gameState;
    var maxPlayers;
    var initialChips;

    const setLocalStorage = (save = 'no') => {
        localStorage.setItem('lobbyId', lobbyId);
        localStorage.setItem('myName', myName);
        localStorage.setItem('state', state);
        localStorage.setItem('errorMessage', errorMessage);
        localStorage.setItem('gameState', JSON.stringify(gameState));
        localStorage.setItem('saveState', save);
    }

    const getLocalStorage = () => {
        if (localStorage.getItem('saveState') === 'yes') {
            lobbyId = localStorage.getItem('lobbyId');
            myName = localStorage.getItem('myName');
            state = localStorage.getItem('state');
            errorMessage = localStorage.getItem('errorMessage');
            gameState = localStorage.getItem('gameState') !== 'undefined' ? JSON.parse(localStorage.getItem('gameState')) : {};
        } else {
            lobbyId = '';
            myName = '';
            state = 'hostjoin';
            errorMessage = '';
            gameState = {};
        }
    }

    const sendAction = (action, amount = 0) => {
        sendMessage(JSON.stringify({
            lobby: lobbyId.toUpperCase(),
            player: {
                name: myName,
            },
            action,
            amount
        }));
    }

    const handleMessage = (m) => {
        if (m.length > 1) {
            let res = JSON.parse(m);
            let {message, game_state, player, error} = res;
            console.log(res);
            if (error && error.player.name == myName || error && error.player.name == "") {
                console.error('Error from server: ' + error.cause + ' --- ' + error.message);
                errorMessage = error.message;
                return;
            }
            gameState = game_state;
            allPlayers = gameState && gameState.players && Object.keys(gameState.players) || allPlayers;
            if (message === 'added player ' + myName) {
                state = 'lobby';
            }
            if (gameState && gameState.playing) {
                state = 'playing';
            }
            setLocalStorage('yes');
        }
        messageStore.set('');
    }

    onMount(() => {
        messageStore.subscribe(currentMessage => {
            handleMessage(currentMessage);
        })
        getLocalStorage();
        socketState.subscribe(state => {
            state === 'connected' && lobbyId !== '' && sendAction('get_state');
        })
    })

    const onJoinPressed = async () => {
        if (myName.length > 0) {
            sendAction('new_player');
        } else {
            const res = await fetch(`http://${serverHost}/${lobbyId}`);
            if (res.status === 404) {
                state = 'hostjoin';
                errorMessage = await res.json()
            } else {
                state = 'entername';
                errorMessage = '';
            }
        }
    }

    const startGame = () => {
        sendAction('start_game');
    }

    const setMaxPlayers = () => {
        sendAction('set_max_players', Number(maxPlayers));
    }

    const setInitialChips = () => {
        sendAction('set_initial_chips', Number(initialChips));
    }

    const onHostPressed = async () => {
        const res = await fetch(`http://${serverHost}/new_game`);
        const json = await res.json();
        lobbyId = JSON.stringify(json).substring(1, 4).toLowerCase();
        state = 'entername';
        errorMessage = '';
    }

    const onHomePressed = () => {
        if (myName !== '') {
            sendMessage(JSON.stringify({
                player: {
                    name: myName
                },
                action: 'remove_player',
                lobby: lobbyId.toUpperCase()
            }));
        }
        resetState();
    }

    const resetState = () => {
        state = 'hostjoin';
        errorMessage = '';
        myName = '';
        lobbyId = '';
        allPlayers = [];
        gameState = {};
        setLocalStorage('no');
    }
</script>

<div class="game">

    <Error bind:text={errorMessage}/>

{#if state === 'hostjoin'}

    <Input helperText="Table Code" bind:val={lobbyId} className="code" onEnter={onJoinPressed}/>
    <Button text="Join" onClick={onJoinPressed} variant="primary-outlined"/>
    <p>- or -</p>
    <Button text="Host New Table" onClick={onHostPressed} variant="primary-filled full-width"/>

{:else if state === 'entername'}

    <h3>{lobbyId && `Table Code: ${lobbyId.toUpperCase()}`}</h3>
    <h3>Enter name:</h3>
    <Input helperText="Name" bind:val={myName} onEnter={onJoinPressed}/>
    <Button text="Join" onClick={onJoinPressed} variant="primary-filled" />

{:else if state === 'lobby'}

    <Button text="Start Game" onClick={startGame} disabled={gameState && !gameState.players[myName].is_host}/>
    <h3>{lobbyId && `Table Code: ${lobbyId.toUpperCase()}`}</h3>
    <Players playerNames={allPlayers} />
    {#if gameState && gameState.players && gameState.players[myName].is_host}
        <Settings {setMaxPlayers} {setInitialChips} bind:maxPlayers bind:initialChips/>
    {/if}

{:else if state === 'playing'}

    <h3>Game</h3>

{/if}
</div>

{#if state !== 'hostjoin'}

<div class="home-button">
    <Button text={state === 'playing' ? 'Leave Game' : 'Leave Lobby'} onClick={onHomePressed} variant="primary-filled"/>
</div>

{/if}



<style>
    .game {
        border-radius: 10px;
        border: 4px solid var(--background);
        background-color: var(--secondary);
        padding: 1em;
        margin-left: 1em;
        margin-right: 1em;
        text-align: center;
        align-items: center;
        display: grid 4em;
    }
    .home-button {
        align-self: flex-start;
        padding: 1em;
    }
</style>