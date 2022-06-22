# Start Xcode simulator

[![Step changelog](https://shields.io/github/v/release/bitrise-steplib/bitrise-step-xcode-start-simulator?include_prereleases&label=changelog&color=blueviolet)](https://github.com/bitrise-steplib/bitrise-step-xcode-start-simulator/releases)

Starts an Xcode simulator.


<details>
<summary>Description</summary>

Starts an Xcode simulator.

</details>

## 🧩 Get started

Add this step directly to your workflow in the [Bitrise Workflow Editor](https://devcenter.bitrise.io/steps-and-workflows/steps-and-workflows-index/).

You can also run this step directly with [Bitrise CLI](https://github.com/bitrise-io/bitrise).

### Examples

Boot simulator in the background and use it in the xcode-test Step:
```yaml
- xcode-start-simulator:
    inputs:
    - destination: platform=iOS Simulator,name=iPhone 8,OS=latest
- xcode-test:
    inputs:
    - project_path: ./ios-sample/ios-sample.xcodeproj
    - scheme: ios-sample
    # Simulator
    - destination: $BITRISE_XCODE_DESTINATION # Use the same destination as the xcode-start-simulator Step
```

Detect if simulator timed out and restart the build:
```yaml
- xcode-start-simulator:
    inputs:
    - destination: platform=iOS Simulator,name=iPhone 8,OS=latest
- trigger-bitrise-workflow:
    is_always_run: true
    run_if: '{{enveq "BITRISE_IS_SIMULATOR_TIMEOUT" "true"}}'
    inputs:
    - api_token: $RESTART_TRIGGER_TOKEN
    - workflow_id: wf
```

## ⚙️ Configuration

<details>
<summary>Inputs</summary>

| Key | Description | Flags | Default |
| --- | --- | --- | --- |
| `destination` | Destination specifier describes the simulator device to be started.  The input value uses the same format as xcodebuild's `-destination` option. | required | `platform=iOS Simulator,name=iPhone 8 Plus,OS=latest` |
| `wait_for_boot_timeout` | When larger than 0, will wait for the simulator boot to complete.  Setting to larger than 0 makes it possible to detect hangs or timeouts when booting simulator. If a timeout occurs, the `BITRISE_IS_SIMULATOR_TIMEOUT` output will be set to true.  Using `0` (the default) enables the Simulator boot to occur in parallel to other Steps. | required | `0` |
| `verbose_log` | If this input is set, the Step will print additional logs for debugging. | required | `no` |
| `reset` | If enabled, will shutdown and erase a simulator's contents and settings.  This option is not needed when starting from a clean state on a CI build. It may be used when running testing multiple apps on the same simulator or for making sure that the simulator is indeed in a clean state when an app fails to install due to an unexpected issue.  When enabled erasing contents takes about a second. | required | `no` |
</details>

<details>
<summary>Outputs</summary>

| Environment Variable | Description |
| --- | --- |
| `BITRISE_IS_SIMULATOR_TIMEOUT` | Set to true/false based on starting Xcode Simulator failed with an unrecoverable error.  It can be used to trigger a new build conditionally:  is_always_run: true run_if: '{{enveq "BITRISE_IS_SIMULATOR_TIMEOUT" "true"}}'  |
| `BITRISE_XCODE_DESTINATION` | Device destination specifier  The destination specifer provided in the `destination` Input. It can be used as Input of other Steps, to avoid duplication. |
</details>

## 🙋 Contributing

We welcome [pull requests](https://github.com/bitrise-steplib/bitrise-step-xcode-start-simulator/pulls) and [issues](https://github.com/bitrise-steplib/bitrise-step-xcode-start-simulator/issues) against this repository.

For pull requests, work on your changes in a forked repository and use the Bitrise CLI to [run step tests locally](https://devcenter.bitrise.io/bitrise-cli/run-your-first-build/).

Learn more about developing steps:

- [Create your own step](https://devcenter.bitrise.io/contributors/create-your-own-step/)
- [Testing your Step](https://devcenter.bitrise.io/contributors/testing-and-versioning-your-steps/)
