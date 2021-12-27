<script>
    import Input from "./input.svelte";
    export let setMaxPlayers = () => {}
    export let setInitialChips = () => {}
    export let maxPlayers = '4';
    export let initialChips = '100';
</script>

<div>
    <hr>
    <h2>Settings:</h2>
    <div class="settings">
        <p>Total players:</p>
        <Input className="settings" onEnter={setMaxPlayers} bind:val={maxPlayers}/>
        <p>Starting chips:</p>
        <Input className="settings" onEnter={setInitialChips} bind:val={initialChips}/>
    </div>
</div>

<style>
    .settings {
        padding-top: 0.5em;
        display: grid;
        grid-template-columns: 1fr 1fr;
    }
</style>