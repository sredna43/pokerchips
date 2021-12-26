<script>
    export let disabled = false;
    export let text = "";
    export let onClick = () => {
        console.log("default on click");
    };
    export let variant = 'primary-filled';
</script>

<button class={variant} on:click={onClick} disabled={disabled} title={disabled ? 'Only host can start a game' : 'Start game'}>{text}</button>

<style>
    button {
        transition-duration: 0.4s;
        border: none;
        border-radius: 10px;
        margin: 2px 2px;
        font-size: 2em;
        text-align: center;
        padding: 0.2em 0.4em;
    }

    .primary-filled {
        color: white;
        background-color: var(--primary);
        border: 2px solid var(--primary)
    }

    .primary-filled:hover:not([disabled]) {
        color: black;
        background-color: white;
    }

    .primary-outlined {
        color: black;
        background-color: white;
        border: 2px solid var(--primary);
    }

    .primary-outlined:hover:not([disabled]) {
        color: white;
        background-color: var(--primary);
    }

    .full-width {
        width: 100%;
    }

    button:disabled {
        opacity: 30%;
    }

</style>