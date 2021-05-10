// Package containerc provides commands to build container images.
// Provide:
//  - pipeline to build and deploy images parallel
//  - (important for security) isolate new images into local repository. Prevent overwrite popular images (like ubuntu, alpine...).
//  - (important for security) isolate login to remote repositories. Define credentials per project. No share login session betwen projects.
package containerc
