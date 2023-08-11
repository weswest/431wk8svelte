<script>
  import {Greet} from '../wailsjs/go/main/App.js'
  import {TestVale} from '../wailsjs/go/main/App.js'
  import {CheckWithVale} from '../wailsjs/go/main/App.js'

  let inputText = "";
  let outputText = "";
  let dataIs = "plural"; // default value

  function evaluateText() {
    console.log("evaluateText function called with input:", inputText);
    // You can now pass the dataIs variable to the CheckWithVale function if needed
    CheckWithVale(inputText, dataIs).then(result => {
        console.log("Received result from backend:", result);
        outputText = result;
    });
  }
</script>

<style>
  .container {
      padding: 20px;
      display: flex;
      flex-direction: column;
      align-items: center;
  }

  textarea {
      width: 90%;
      height: 200px;
      margin-bottom: 20px;
  }

  button {
      margin-bottom: 20px;
  }

  .output {
      width: 90%;
      border: 1px solid #ccc;
      padding: 10px;
      white-space: pre-line;
  }
</style>

<div class="container">
  <div>
      <p>Static text that describes how to use the application.</p>
      <label>
        <input type="radio" bind:group={dataIs} value="singular">
        'Data' is Singular
      </label>
      <label>
        <input type="radio" bind:group={dataIs} value="plural" checked>
        'Data' is Plural
      </label>
      <label>
        <input type="radio" bind:group={dataIs} value="all">
        Produce All Warnings
      </label>
  </div>
  <textarea bind:value={inputText} placeholder="Enter your test here..."></textarea>
  <button on:click={evaluateText}>Evaluate</button>
  <div class="output">{outputText}</div>
</div>
