<script>
    import Game from "../lib/components/game.svelte";
</script>

<main>
    <h1>Online Poker Chips</h1>
    <Game/>
</main>

<style>
    :global(:root) {
        --primary: #3C4F76;
        --secondary: #DDDBF1;
        --header-text: #FF521B;
        --background: #383F51;
        --accent: #2EC0F9;
        font-family: 'Lucida Sans', 'Lucida Sans Regular', 'Lucida Grande', 'Lucida Sans Unicode', Geneva, Verdana, sans-serif;
    }
    main {
        max-width: 960px;
        margin: 1em auto;
        padding: 1em;
    }
    h1 {
        padding-top: 0.5em;
        color: var(--header-text);
        text-transform: uppercase;
        font-size: 4em;
        font-weight: 100;
        text-align: center;
    }
</style>