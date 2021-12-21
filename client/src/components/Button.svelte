<script lang="ts">
    export let text: string;
    export let classname: string = "";
    export let onClick = () => {console.log("default click")};
</script>

<button class={classname} on:click={onClick}>
    {text}
</button>

<style>
    button {
        border: none;
        padding: 16px 32px;
        text-align: center;
        text-decoration: none;
        display: inline-block;
        font-size: 16px;
        margin: 4px 2px;
        transition-duration: 0.4s;
        cursor: pointer;
    }

    .filled {
        background-color: cadetblue;
        color: white;
        border: 2px solid cadetblue;
    }

    .outlined {
        background-color: white;
        color: black;
        border: 2px solid cadetblue;
    }

    .filled:hover {
        background-color: white;
        color: black;
    }

    .outlined:hover {
        background-color: cadetblue;
        color: white;
    }


</style>