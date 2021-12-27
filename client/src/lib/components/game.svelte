<script>
    import Button from './button.svelte';
    import Input from './input.svelte';
    import Error from './error.svelte';
    import Players from './players.svelte';
    import Settings from './settings.svelte';
    import Actions from './actions.svelte';
    import { onMount } from 'svelte';
    import { serverHost, sendMessage, messageStore, socketState } from '../store.js';

    var lobbyId = '';
    var myName = '';
    var state = 'hostjoin';
    var errorMessage = '';
    var gameState = {};
    var maxPlayers = '4';
    var initialChips = '100';
    var betText = 'Bet';
    var canCheck = true;

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
            let {lobby, message, game_state, player, error} = res;
            if (lobby.toLowerCase() !== lobbyId) {
                return;
            }
            if (error && (error.player.name === myName && !error.message.includes('exists') || (error.message.includes('exists') && state == 'entername'))) {
                console.error('Error from server: ' + error.cause + ' --- ' + error.message);
                errorMessage = error.message;
                return;
            }
            if (game_state !== null) {
                gameState = game_state;
            }
            if (gameState && gameState.playing && state !== 'playing') {
                state = 'playing';
            }
            if (message === 'added player ' + myName) {
                state = 'lobby';
            }
            if (message.includes('removed') && state !== 'hostjoin') {
                sendAction('get_state');
            }
            if (message.includes('removed') && player.name == myName) {
                resetState();
            }
            if (message.includes('bet')) {
                canCheck = false;
                betText = 'Raise';
            }
            if (message.includes('check')) {
                betText = 'Bet';
            }
        }
        messageStore.set('');
    }

    onMount(() => {
        messageStore.subscribe(currentMessage => {
            handleMessage(currentMessage)
        });
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
        gameState = {};
    }

    const isMyTurn = () => {
        return gameState && gameState.players && gameState.players[myName].spot === gameState.whose_turn;
    }

    const amDealer = () => {
        return gameState && gameState.players && gameState.players[myName].spot === gameState.dealer;
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
    <Players bind:gameState={gameState} />
    {#if gameState && gameState.players && gameState.players[myName].is_host}
        <Settings {setMaxPlayers} {setInitialChips} bind:maxPlayers bind:initialChips/>
    {/if}

{:else if state === 'playing'}

    <h2>Pot: {gameState.pot}</h2>
    <Players bind:gameState={gameState}/>
    <Actions bind:betText bind:canCheck/>

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