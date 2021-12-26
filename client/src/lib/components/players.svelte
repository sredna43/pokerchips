<script>
    export let playerNames = [];
</script>

<h4>Players:</h4>
<ol>
    {#each playerNames as player }
        <li>{player}</li>
    {/each}
</ol>