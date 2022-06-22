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
    - wait_for_boot_timeout: 90
- trigger-bitrise-workflow:
    is_always_run: true
    run_if: '{{enveq "BITRISE_SIMULATOR_STATUS" "hanged"}}'
    inputs:
    - api_token: $RESTART_TRIGGER_TOKEN
    - workflow_id: wf
```