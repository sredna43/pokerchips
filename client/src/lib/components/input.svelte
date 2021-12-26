<script>
    export let val = '';
    export let helperText = '';
    export let className = '';
    export let onEnter = () => {
        console.log("Enter key pressed")
    }
    const keyUp = (e) => {
        e.keyCode === 13 && onEnter()
    }
</script>

<input bind:value={val} placeholder={helperText} class={className} on:keyup={keyUp}>

<style>
    input {
        color: black;
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
    .settings {
        font-size: 2em;
        padding: 0% 0%;
        width: 40%;
    }
</style>