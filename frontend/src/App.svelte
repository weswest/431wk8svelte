<script>
  import {Greet} from '../wailsjs/go/main/App.js'
  import {TestVale} from '../wailsjs/go/main/App.js'
  import {CheckWithVale} from '../wailsjs/go/main/App.js'

  let inputText = "";
  let outputText = "";
  let dataIs = "plural"; // default value

  const sample1 = `In today's digital age, data is crucial. Data drives most decisions, from business strategies to personal choices. However, there's often confusion about how to treat the word "data" in written English.

Traditionally, "data" is the plural form of "datum." This means that when referring to multiple pieces of information, one would say, "The data are clear." On the other hand, if referring to a single piece of information, the correct usage would be, "This datum is clear."

However, modern usage has evolved. Many people now treat "data" as a singular noun, especially in the context of digital information. For example, "The data is being processed" or "This data shows a significant trend."

Google's style guide provides clarity on many such issues. It's a comprehensive resource that offers guidance on grammar, punctuation, and word choice. Adhering to a style guide ensures consistency and clarity in written communication.

This app respects individual preferences. Users can choose whether they want to treat "data" as a singular or plural noun. This customization ensures that documents align with the writer's intent.

For instance, consider the sentence, "The data is compelling." If a user prefers the traditional approach, the tool flags this as an error. Conversely, if the preference is for modern usage, the sentence "The data show a rising trend" would be highlighted.

In conclusion, language is dynamic, and its usage evolves over time. While it's essential to respect traditional rules, it's equally important to acknowledge contemporary practices. The tool aims to bridge this gap, offering flexibility while maintaining adherence to the Google style guide.

This text triggers the singular or plural rule for "data" based on the user's preference. If the user chooses the traditional approach (data as plural), the sentences "The data is being processed" and "The data is compelling" is flagged. If the user opts for the modern usage (data as singular), the sentence "The data show a rising trend" is highlighted.`; 
  const sample2 = `In a world where numbers never lie,
The data are a beacon, shining bright.
Once data were obscured, now they fly,
Revealing truths, bringing them to light.

Data have secrets, whispered in the night,
Yet data themselves can sometimes deceive.
In the dance of figures, data do unite,
Telling tales that few can believe.

But data is more than just cold hard facts,
For once data was a story untold.
Now, data has a voice, and it acts,
With power, data itself does hold.

In the end, data does what it must,
To earn our ever-growing trust.`; 

  function evaluateText() {
    console.log("evaluateText function called with input:", inputText);
    CheckWithVale(inputText, dataIs).then(result => {
        console.log("Received result from backend:", result);
        outputText = result;
    });
  }

  function loadSample(sample) {
    inputText = sample;
  }

  function loadFile(event) {
    const file = event.target.files[0];
    if (file) {
        const reader = new FileReader();
        reader.onload = function(e) {
            if (typeof e.target.result === 'string') {
                inputText = e.target.result;
            } else {
                inputText = "File's text not found";
            }
        };
        reader.readAsText(file);
    }
}

</script>

<style>
  .container {
      padding: 20px;
      display: flex;
      flex-direction: column;
      align-items: center;
  }

  .input-section {
      display: flex;
      align-items: start;
      width: 90%;
  }

  textarea {
      flex: 1;
      height: 200px;
      margin-right: 20px;
  }

  .sample-buttons {
      display: flex;
      flex-direction: column;
  }

  button {
      margin-bottom: 20px;
  }

  .output {
      width: 90%;
      border: 1px solid #ccc;
      padding: 10px;
      white-space: pre-line;
      text-align: left; 
  }

  .rules {
      text-align: left;
      font-size: 0.8rem;
  }
</style>

<div class="container">
  <div>
      <h2>Assignment:</h2>
      <p>
        Create a program using Wails, Svelte, and Vale that allows users to enter text and receive feedback on their grammar. As a test case, focus is on whether "data" should be plural or singular.
      <h2>How to Use</h2>
      <p>
        Enter text and click evaluate. The only rules turned on are the data singular / plural, as defined by the radio boxes below. Please see the bottom for the specific rules that are captured.  Note that these rules are hard-coded: with more time we can add more sophistication.
      </p>
      <p>
        If you want some sample text to get you started, please click the samples on the right.
      </p>
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
  <div>
      <button>
        <input type="file" id="fileInput" on:change={loadFile} accept=".txt,.md" style="width: 100%;">
        Read in .txt or .md
      </button>
  </div>
  <div class="input-section">
      <textarea bind:value={inputText} placeholder="Enter your text here..."></textarea>
      <div class="sample-buttons">
          <button on:click={() => loadSample(sample1)}>Sample 1</button>
          <button on:click={() => loadSample(sample2)}>Sample 2</button>

      </div>
  </div>
  <button on:click={evaluateText}>Evaluate</button>
  <div class="output">{outputText}</div>
  <div class="rules">
      <h3>This program adheres to the following rules:</h3>
      <ul>
          <li>Data is / are</li>
          <li>Data was / were</li>
          <li>Data has / have</li>
          <li>Data itself / themselves</li>
          <li>Data does / do</li>
      </ul>
  </div>
</div>
