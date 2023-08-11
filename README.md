# Longevity Digital Twin

This repository stores various Longevity Digital Twins. It is part of a master-thesis with the goal to develop an Orchestration and Deployment Manager.

For a complete understanding of the project, please refer to the following repositories:

- [Orchestration and Deployment Manager](https://github.com/pixelboehm/longevity): Main application that handles the orchestration of LDTs.
- [LDT Meta Repository](https://github.com/pixelboehm/meta-ldt): Stores a file with links to repositories containing LDTs.
- [ESP32 Applications](https://github.com/pixelboehm/longevity-esp32): Stores ESP32 applications that are our smart devices.
- (Optional) [Homebrew-LDT](https://github.com/pixelboehm/homebrew-ldt): Contains Homebrew (outdated) formulas for the ODM and LDTs. The formulas are not up-to-date anymore, but can be enabled through the `.goreleaser.yml` again.

## LDTs

- Lightbulb (HomeKit Accessory)
- Switch (HomeKit Accessory)

## Dependencies

- goreleaser (for local builds only)
- make

## Building

Pushing / Merging into the main branch triggers a GitHub workflow that automatically builds and releases a new version of the LDTs of this repository.

As every LDT requires a Web-Of-Things description, regular building via `go` is possible, but not advised. For a local build is possible via `make releaseLocal`.

## Compatibility

All releases with `0.x.x` are considered unstable. This is due to the continuous changes in the architecture of LDTs and the ODM. Therefore, it is advised to use the latest version of LDTs, if the ODM is also the latest version.
As not all requirements are fulfilled at this point (according to the evaluation of the thesis), I would not release a version `1.0.0` of the LDTs.