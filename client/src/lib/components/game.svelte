<script>
    import Button from './button.svelte';
    import { onMount } from 'svelte';
    import { serverHost, sendMessage, messageStore } from '../store.js';

    var lobbyId = '';
    var name = '';
    var state = 'hostjoin';
    var error = '';

    const handleMessage = (m) => {
        if (m.length > 1 ) {
            let res = JSON.parse(m);
            let {message, game_state, player, error} = res;
            console.log('message received: ' + message);
            if (error) {
                console.error('Error from server: ' + error.cause + ' --- ' + error.message);
                return;
            }
            if (message) {
                console.log(message)
            }
        }
        messageStore.set('');
    }

    onMount(() => {
        messageStore.subscribe(currentMessage => {
            handleMessage(currentMessage);
        })
    })

    const onJoinPressed = async () => {
        if (name.length > 0) {
            state = 'playing';
        } else {
            const res = await fetch(`http://${serverHost}/${lobbyId}`);
            if (res.status === 404) {
                state = 'hostjoin';
            } else {
                state = 'entername';
            }
        }
    }

    const onHostPressed = async () => {
        const res = await fetch(`http://${serverHost}/new_game`);
        const json = await res.json();
        lobbyId = JSON.stringify(json).substring(1, 4);
        state = 'playing';
    }

    const onHomePressed = () => {
        state = 'hostjoin';
        lobbyId = '';
    }
</script>

<div class="game">
{#if state === 'hostjoin'}
    <input placeholder="Table Code" bind:value={lobbyId} class="code">
    <Button text="Join" onClick={onJoinPressed} variant="primary-outlined"/>
    <p>- or -</p>
    <Button text="Host New Table" onClick={onHostPressed} variant="primary-filled full-width"/>
{:else if state === 'entername'}
    <h3>{lobbyId && `Table Code: ${lobbyId.toUpperCase()}`}</h3>
    <h3>Enter name:</h3>
    <input placeholder="Name" bind:value={name}>
    <Button text="Join" onClick={onJoinPressed} variant="primary-filled" />
{:else}
    <h3>{lobbyId && `Table Code: ${lobbyId.toUpperCase()}`}</h3>
{/if}
</div>

{#if state !== 'hostjoin'}
<div class="back-button">
    <Button text="Home" onClick={onHomePressed} variant="primary-filled"/>
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
    }
    .back-button {
        align-self: flex-start;
        padding: 1em;
    }
    input {
        color: black;
        width: auto;
        border: 2px solid var(--primary);
        border-radius: 5px;
        text-align: center;
        font-size: 2em;
        text-align: center;
        padding: 0.2em 0em;
    }
    .code {
        text-transform: uppercase;
    }
</style>