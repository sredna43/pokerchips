<script>
    export let gameState = {};
</script>

<h4>Players:</h4>
<ol>
    {#each Object.keys(gameState?.players) as player }
        <li>{player} - Chips: {gameState && gameState.players && gameState.players[player].chips || 0} {gameState && gameState.players && gameState.players[player].folded ? 'Folded' : ''}</li>
    {/each}
</ol>