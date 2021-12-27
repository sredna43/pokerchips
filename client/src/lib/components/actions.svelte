<script>
    import Button from './button.svelte';
    export let check = () => {}
    export let bet = () => {}
    export let fold = () => {}
    export let canCheck = true;
    export let betText = 'Bet';
</script>

<div>
    <Button text={canCheck ? "Check" : "Call"} onClick={canCheck ? check : bet}/>
    <Button text={betText} onClick={bet}/>
    <Button text = "Fold" onClick={fold}/>
</div>

<style>
    div {
        padding-top: 1em;
    }
</style>