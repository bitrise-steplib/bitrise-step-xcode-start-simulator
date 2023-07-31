### Examples

Boot simulator in the background and use it in the xcode-test Step:
```yaml
- xcode-start-simulator:
    inputs:
    - destination: platform=iOS Simulator,name=Bitrise iOS default,OS=latest
- xcode-test:
    inputs:
    - project_path: ./ios-sample/ios-sample.xcodeproj
    - scheme: ios-sample
    - destination: $BITRISE_XCODE_DESTINATION # Use the same destination as the xcode-start-simulator Step
```

Boot Rosetta Simulator and use it in the xcode-test Step:
```yaml
- xcode-start-simulator:
    inputs:
    - destination: platform=iOS Simulator,name=Bitrise iOS default,OS=latest,arch=x86_64
- xcode-test:
    inputs:
    - project_path: ./ios-sample/ios-sample.xcodeproj
    - scheme: ios-sample
    - destination: $BITRISE_XCODE_DESTINATION # Use the same destination as the xcode-start-simulator Step
    # Disabling parallel testing ensures that prebooted device is used. ARCHS=x86_64 is optional, to enable project compilation
    - xcodebuild_options: -verbose -parallel-testing-enabled NO  ARCHS=x86_64
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