package mail

var accountVerificationBody = `
<!DOCTYPE html>
<html>
  <head>
    <style>
      body {background-color: #fff; font-size: 18px; line-height: 1.5; color: #24292e;}
      div.content {padding: 10px; margin: auto; max-width: 800px; }
      div.amp-verify {padding: 8px; margin-bottom: 10px; border-style: solid; border-width: 1px; border-color: #c0d3eb;}
      div.amp-documentation {padding: 8px; border-style: solid; border-width: 1px; border-color: #c0d3eb;}
      p.welcome-message {background-color: #4e7cad; color:#fff; font-size: 150%; text-align: center; border-width: 2px; border-color: #c0d3eb; padding: 10px 20px; font-family: 'Montserrat', sans-serif;}
      footer {display: block; background-color: #4e7cad; color: #fff; padding: 20px; font-family: 'Montserrat', sans-serif;}
      footer p {text-align: center; margin: 0 auto; padding: 0 20px;}
      footer p a {color: #fff;}
      pre {display: block; background-color: #f6f8fa; padding: 16px; margin-top: 0; margin-bottom: 16px; overflow: auto; white-space: pre; font-size: 85%; line-height: 1.45; box-sizing: border-box; border-radius: 3px; word-wrap: normal;}
      code {font-family: "SFMono-Regular", Consolas, "Liberation Mono", Menlo, Courier, monospace; display: inline; padding: 0; margin: 0; overflow: visible; line-height: inherit; word-break: break-all; background-color: transparent; border: 0;}
    </style>
  </head>
  <body>
    <div class="content">
        <p class="welcome-message">Thanks for joining AMP</p>
      <div class="amp-verify">
        <p>Your account <i>{accountName}</i> has been successfully created.</p>
        <p>One last step before you can use it, this is to make sure that your email address is valid. The command below will verify your account.</p>
        <pre><code>amp -s {cliAddr} user verify {token}</code></pre>
        <p>If you didn't make this request, you can safely ignore this email.</p>
        <p>This token is valid for one hour only.</p>
      </div>
      <div class="amp-documentation">
        <p>The latest CLI can be downloaded at <a src="https://github.com/appcelerator/amp/releases">the official amp github repository</a>. Running <code>amp -h</code> will display extensive information on how to use it.</p>
        <p>Documentation on the amp platform is available at <a src="http://appcelerator.io">appcelerator.io</a>.</p>
      </div>
      <footer>
      <p>AMP is an open source project sponsored by <a href="https://www.axway.com">Axway</a>. It is licensed under <a href="https://github.com/appcelerator/amp/blob/master/LICENSE">Apache 2.0</a>.</p>
      </footer>
    </div>
  </body>
</html>
`
