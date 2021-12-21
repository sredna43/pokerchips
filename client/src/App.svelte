<script lang="ts">
	import { onMount } from 'svelte';
	import store from './store.js'
	import Button from './components/Button.svelte'
	let header: string = "online poker chips"

	const onHostPressed = () => {
		console.log("host pressed");
		store.sendMessage('{"action": "new_game"}')
	};

	const onJoinPressed = () => {
		console.log("join pressed");
	};
</script>

<main>
	<h1>{header}</h1>
	<Button text="Join" classname="filled" onClick={onJoinPressed}/>
	<Button text="Host" classname="outlined" onClick={onHostPressed}/>
</main>

<style>
	main {
		text-align: center;
		padding: 1em;
		max-width: 240px;
		margin: 0 auto;
		background-color: white;
	}

	h1 {
		color: #ff3e00;
		text-transform: uppercase;
		font-size: 4em;
		font-weight: 100;
	}

	@media (min-width: 640px) {
		main {
			max-width: none;
		}
	}
</style>