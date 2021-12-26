<script>
    export let gameState = {};
    let playerNames = Object.keys(gameState?.players) || []
    console.log(gameState);
</script>

<h4>Players:</h4>
<ol>
    {#each playerNames as player }
        <li>{player} - Chips: {gameState && gameState.players && gameState.players[player].chips || 0} {gameState && gameState.players && gameState.players[player].folded ? 'Folded' : ''}</li>
    {/each}
</ol>