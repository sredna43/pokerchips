<script>
    import Button from './button.svelte';
    export let check = () => {}
    export let bet = () => {}
    export let fold = () => {}
    export let canCheck = true;
</script>

<div>
    {#if canCheck}
        <Button text="Check" onClick={check}/>
    {/if}
    <Button text="Bet" onClick={bet}/>
    <Button text = "Fold" onClick={fold}/>
</div>

<style>
    div {
        padding-top: 1em;
    }
</style>