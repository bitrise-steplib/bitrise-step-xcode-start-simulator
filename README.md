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
    run_if: '{{enveq "BITRISE_IS_SIMULAOR_TIMEOUT" "true"}}'
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
| `erase` | If enabled, will erase a simulator's contents and settings. | required | `no` |
| `wait_for_boot` | If set, will wait simulator boot to complete.  Setting to `yes` makes possible to detect hangs or timeouts when booting simulator. If a timeout occurs, the `BITRISE_IS_SIMULATOR_TIMEOUT` output will be set to true.  Using `no` (the default) enables to boot the Simulator parallel to other Steps. | required | `no` |
| `verbose_log` | If this input is set, the Step will print additional logs for debugging. | required | `no` |
| `wait_for_boot_timeout` | Maximum allowed time for simulator boot (in seconds)  "Wait for simulator to boot" (`wait_for_boot`) must be set to `yes`. | required | `90` |
</details>

<details>
<summary>Outputs</summary>

| Environment Variable | Description |
| --- | --- |
| `BITRISE_IS_SIMULATOR_TIMEOUT` | Set to true/false based on starting Xcode Simulator failed with an unrecoverable error.  |
| `BITRISE_XCODE_DESTINATION` | Device destination specifier  The destination specifer provided in the `destination` Input, so it can be used in other Steps too. |
</details>

## 🙋 Contributing

We welcome [pull requests](https://github.com/bitrise-steplib/bitrise-step-xcode-start-simulator/pulls) and [issues](https://github.com/bitrise-steplib/bitrise-step-xcode-start-simulator/issues) against this repository.

For pull requests, work on your changes in a forked repository and use the Bitrise CLI to [run step tests locally](https://devcenter.bitrise.io/bitrise-cli/run-your-first-build/).

Learn more about developing steps:

- [Create your own step](https://devcenter.bitrise.io/contributors/create-your-own-step/)
- [Testing your Step](https://devcenter.bitrise.io/contributors/testing-and-versioning-your-steps/)
