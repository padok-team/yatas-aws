# Changelog

All notable changes to this project will be documented in this file. See [standard-version](https://github.com/conventional-changelog/standard-version) for commit guidelines.

## [1.13.0](https://github.com/padok-team/yatas-aws/compare/v1.12.0...v1.13.0) (2025-06-27)


### Features

* **s3:** check bucket deletion policy ([#192](https://github.com/padok-team/yatas-aws/issues/192)) ([ef027cb](https://github.com/padok-team/yatas-aws/commit/ef027cbc5bf841f536615dd84d9fe1ace2524067))

## [1.12.0](https://github.com/padok-team/yatas-aws/compare/v1.11.1...v1.12.0) (2025-05-09)


### Features

* **db-backups:** add 3 checks for 3 backup encryption ([#203](https://github.com/padok-team/yatas-aws/issues/203)) ([f4eb252](https://github.com/padok-team/yatas-aws/commit/f4eb252fb15105a2492ce21b8b91f48fbe438aef))
* **rds:** add check for audit logs ([#201](https://github.com/padok-team/yatas-aws/issues/201)) ([0fe1ceb](https://github.com/padok-team/yatas-aws/commit/0fe1cebca6e75c6b6744bb3645d549e299709e60))

### [1.11.1](https://github.com/padok-team/yatas-aws/compare/v1.11.0...v1.11.1) (2025-03-14)


### Bug Fixes

* **ssm:** move misplaced blocking check channel call ([#199](https://github.com/padok-team/yatas-aws/issues/199)) ([820d731](https://github.com/padok-team/yatas-aws/commit/820d7315c67bfe4ae2a08536babe449079b3c16b))

## [1.11.0](https://github.com/padok-team/yatas-aws/compare/v1.9.0...v1.11.0) (2025-03-07)


### Features

* **acm:** add check for secure encryption in ACM certificates ([#188](https://github.com/padok-team/yatas-aws/issues/188)) ([7a4238d](https://github.com/padok-team/yatas-aws/commit/7a4238d6dd4b40842369b1a39128ff9ef6ef8a09))
* **cloudtrail:** test if cloudtrails are enabled ([#152](https://github.com/padok-team/yatas-aws/issues/152)) ([a86e662](https://github.com/padok-team/yatas-aws/commit/a86e6629d1de10591dd0e034abac61144362fb8c))
* **elb:** rename loadbalancers to elb & add alb ensure https check ([#190](https://github.com/padok-team/yatas-aws/issues/190)) ([fb2bf99](https://github.com/padok-team/yatas-aws/commit/fb2bf997fb37da8224798a2f9df993c9fc4f86fe))
* **guardduty:** add check for high severity findings ([#187](https://github.com/padok-team/yatas-aws/issues/187)) ([1d389a8](https://github.com/padok-team/yatas-aws/commit/1d389a8f9d84854007917c27eea8bf34e532c36a))
* **ssm:** add bastion audit logs check ([#186](https://github.com/padok-team/yatas-aws/issues/186)) ([f743a98](https://github.com/padok-team/yatas-aws/commit/f743a98aa7cce04311ba0062d9f3aefd802a4c92))


### Bug Fixes

* **AWS_IAM_006:** nil pointer in GetPasswordPolicy ([#189](https://github.com/padok-team/yatas-aws/issues/189)) ([bdf0c3d](https://github.com/padok-team/yatas-aws/commit/bdf0c3deb3cec5932d5a16e275b086baeaf1de26))

## [1.10.0](https://github.com/padok-team/yatas-aws/compare/v1.9.0...v1.10.0) (2025-03-07)


### Features

* **acm:** add check for secure encryption in ACM certificates ([#188](https://github.com/padok-team/yatas-aws/issues/188)) ([7a4238d](https://github.com/padok-team/yatas-aws/commit/7a4238d6dd4b40842369b1a39128ff9ef6ef8a09))
* **cloudtrail:** test if cloudtrails are enabled ([#152](https://github.com/padok-team/yatas-aws/issues/152)) ([a86e662](https://github.com/padok-team/yatas-aws/commit/a86e6629d1de10591dd0e034abac61144362fb8c))
* **guardduty:** add check for high severity findings ([#187](https://github.com/padok-team/yatas-aws/issues/187)) ([1d389a8](https://github.com/padok-team/yatas-aws/commit/1d389a8f9d84854007917c27eea8bf34e532c36a))
* **ssm:** add bastion audit logs check ([#186](https://github.com/padok-team/yatas-aws/issues/186)) ([f743a98](https://github.com/padok-team/yatas-aws/commit/f743a98aa7cce04311ba0062d9f3aefd802a4c92))


### Bug Fixes

* **AWS_IAM_006:** nil pointer in GetPasswordPolicy ([#189](https://github.com/padok-team/yatas-aws/issues/189)) ([bdf0c3d](https://github.com/padok-team/yatas-aws/commit/bdf0c3deb3cec5932d5a16e275b086baeaf1de26))

## [1.9.0](https://github.com/padok-team/yatas-aws/compare/v1.8.1...v1.9.0) (2024-12-18)


### Features

* **configservice:** test if aws config is enabled ([#154](https://github.com/padok-team/yatas-aws/issues/154)) ([6674659](https://github.com/padok-team/yatas-aws/commit/66746595de3945b93ed11920bf65acbcad4270de))
* **iam:** new test on IAM password policy ([#150](https://github.com/padok-team/yatas-aws/issues/150)) ([974557a](https://github.com/padok-team/yatas-aws/commit/974557a8636f58d5711b0260623b6ba22023890e))
* **iam:** new test on no console pwd for non human ([#151](https://github.com/padok-team/yatas-aws/issues/151)) ([b454363](https://github.com/padok-team/yatas-aws/commit/b454363378edaaa1d1a7bcd092d55750db8cf9e3))
* **vpc:** add new test on VPC including private and public subnet ([#149](https://github.com/padok-team/yatas-aws/issues/149)) ([878d632](https://github.com/padok-team/yatas-aws/commit/878d6325de00171332c93a292f9281eb52078fcf))

### [1.8.1](https://github.com/padok-team/yatas-aws/compare/v1.8.0...v1.8.1) (2023-04-11)


### Bug Fixes

* **lambda:** check if not null ([3f70700](https://github.com/padok-team/yatas-aws/commit/3f70700322cb16b3f79d1786669b58320b3cbfc8))

## [1.8.0](https://github.com/padok-team/yatas-aws/compare/v1.7.0...v1.8.0) (2023-04-10)


### Features

* **cognito:** add cognito check for self-registration ([7334bed](https://github.com/padok-team/yatas-aws/commit/7334bed438516cfb1552e973f2ada2c8611b79b3))
* **lambda:** add lambda checks for secrets in environment and URL AuthType ([610606f](https://github.com/padok-team/yatas-aws/commit/610606f2b18d4a453a02dc37f9a0803aedc4fc0a))
* **lambda:** new test ([43d6de3](https://github.com/padok-team/yatas-aws/commit/43d6de31c55105109bfc3e7a87073116fc91263b))
* return ARNs as ResourceID instead of name ([f8dcc44](https://github.com/padok-team/yatas-aws/commit/f8dcc443f812b949061a4a1c91c06d84144edaa0))

## [1.7.0](https://github.com/padok-team/yatas-aws/compare/v1.6.0...v1.7.0) (2023-04-09)


### Features

* **deps:** update ([4af13e1](https://github.com/padok-team/yatas-aws/commit/4af13e1611719b9492e599b83e41ba63acbdb0f6))
* **go:** updated to 1.20 ([03895f1](https://github.com/padok-team/yatas-aws/commit/03895f16b6e8d70a086b6a504f1b4b00ed3a86bd))
* **logger:** update new logger ([45ff0ac](https://github.com/padok-team/yatas-aws/commit/45ff0ac10666e940b42ff4234d2656a7e02e3045))
* **new-yatas:** update imports and function calls from YATAS ([926aebc](https://github.com/padok-team/yatas-aws/commit/926aebccf86e17caba64c5d9ac08b731c5a1a26d))
* **panic:** to logger ([1381d16](https://github.com/padok-team/yatas-aws/commit/1381d16331bde6fcb42f0b2a7d0708598ca3b659))
* **update:** dependcies ([5cecc6e](https://github.com/padok-team/yatas-aws/commit/5cecc6eb061a3d523942f03dcc029a607b4a721f))

## [1.6.0](https://github.com/padok-team/yatas-aws/compare/v1.5.5...v1.6.0) (2023-03-27)


### Features

* **s3:** S3_002: better output and error handling + add tests ([4fcedd4](https://github.com/padok-team/yatas-aws/commit/4fcedd425518d5506f7f19c8613a74c25e24c23b))
* **s3:** S3_002: check bucket has no replication to other region ([a086eaf](https://github.com/padok-team/yatas-aws/commit/a086eaf1d8eda900680684ee951c77e50bae4c18))
* **s3:** S3_002: remove old version of check ([fae1929](https://github.com/padok-team/yatas-aws/commit/fae1929ec7896faa4230c15d14faf164bbf9d1b9))
* **s3:** S3_002: update README ([7ee66e0](https://github.com/padok-team/yatas-aws/commit/7ee66e04db050cdd6bb3731bc76efef5414f7c46))

### [1.5.5](https://github.com/padok-team/yatas-aws/compare/v1.5.4...v1.5.5) (2023-03-24)


### Bug Fixes

* **s3:** nil pointer gets3 ([1fa9b06](https://github.com/padok-team/yatas-aws/commit/1fa9b069d3ecca1b41a39ff3d26930b242071015))

### [1.5.4](https://github.com/padok-team/yatas-aws/compare/v1.5.3...v1.5.4) (2023-03-14)


### Bug Fixes

* **rds:** fixed when no rights to list rds ([3ed2787](https://github.com/padok-team/yatas-aws/commit/3ed2787835902d63fb0f3ee933e567e21dbfead5))

### [1.5.3](https://github.com/padok-team/yatas-aws/compare/v1.5.2...v1.5.3) (2023-02-24)


### Bug Fixes

* **getters:** use getter result only after error handling ([c79e53c](https://github.com/padok-team/yatas-aws/commit/c79e53c83e71131d81ffcbe4fb5e52949a6c7957))

### [1.5.2](https://github.com/padok-team/yatas-aws/compare/v1.5.1...v1.5.2) (2023-02-23)


### Bug Fixes

* **getters:** made all getters fault tolerant by returning empty struct instead of only printing ([7c4336b](https://github.com/padok-team/yatas-aws/commit/7c4336b037d902b81c03c9420af8c4efb505785a))

### [1.5.1](https://github.com/padok-team/yatas-aws/compare/v1.5.0...v1.5.1) (2023-02-06)

## [1.5.0](https://github.com/padok-team/yatas-aws/compare/v1.4.0...v1.5.0) (2023-02-06)


### Features

* **upgrade:** upgraded all services ([e6758c9](https://github.com/padok-team/yatas-aws/commit/e6758c9b034ab0eeb876ab4be24bab8b33b33efc))


### Bug Fixes

* **s3:** fixed issue with access policy ([d8075a9](https://github.com/padok-team/yatas-aws/commit/d8075a9f0a95de05656cb70b8e03a9fc9c93f56a))

## [1.4.0](https://github.com/padok-team/yatas-aws/compare/v1.3.0...v1.4.0) (2023-01-11)


### Features

* **iam:** add check for IAM role privilege escalation ([0140c08](https://github.com/padok-team/yatas-aws/commit/0140c080327073d073dc3a880bf3bff122769efa))
* **iam:** update required permissions for privelege escalation ([43f44a6](https://github.com/padok-team/yatas-aws/commit/43f44a6814139be6755112a63bb7383b08118d39))

## [1.3.0](https://github.com/padok-team/yatas-aws/compare/v1.2.2...v1.3.0) (2022-12-23)


### Features

* **dependencies:** updated ([19578d3](https://github.com/padok-team/yatas-aws/commit/19578d30888ddddde97d2d33ea762a4f23da131b))

### [1.2.2](https://github.com/padok-team/yatas-aws/compare/v1.2.1...v1.2.2) (2022-12-15)

### [1.2.1](https://github.com/padok-team/yatas-aws/compare/v1.2.0...v1.2.1) (2022-10-12)


### Bug Fixes

* **error:** panic to print ([51ecb24](https://github.com/padok-team/yatas-aws/commit/51ecb24c1dda11d8097121916fb1e4fa493efbf2))

## [1.2.0](https://github.com/padok-team/yatas-aws/compare/v1.1.0...v1.2.0) (2022-10-04)


### Features

* **cognito:** added new test for unauthenticated ([ba60f4a](https://github.com/padok-team/yatas-aws/commit/ba60f4a8f0be230c10bb4503125c5122f3a0d77e))

## [1.1.0](https://github.com/padok-team/yatas-aws/compare/v1.0.0...v1.1.0) (2022-09-27)


### Features

* **plugin:** added new feature for categories ([a25d2e5](https://github.com/padok-team/yatas-aws/commit/a25d2e5a2bcb1316af09b80267752367ead6789f))

## [1.0.0](https://github.com/padok-team/yatas-aws/compare/v0.0.8...v1.0.0) (2022-09-26)


### Features

* **plugin:** updated to latest version ([bbed769](https://github.com/padok-team/yatas-aws/commit/bbed7695dbe636a3a63685f3d91a7e5aef05dd52))

### [0.0.8](https://github.com/padok-team/yatas-aws/compare/v0.0.7...v0.0.8) (2022-09-26)


### Features

* **plugins:** upgraded to new interface ([c0db585](https://github.com/padok-team/yatas-aws/commit/c0db58594f2b7e2b6205f41cdffdd945ceac6c6c))

### [0.0.7](https://github.com/padok-team/yatas-aws/compare/v0.0.6...v0.0.7) (2022-09-09)


### Features

* **auth:** plugin can now get authentification from interface passed ([7b006f3](https://github.com/padok-team/yatas-aws/commit/7b006f395c9e100a68a59a0f6fd4cae056a94228))

### [0.0.6](https://github.com/padok-team/yatas-aws/compare/v0.0.5...v0.0.6) (2022-09-09)


### Features

* **makefile:** build in good destination ([f15217a](https://github.com/padok-team/yatas-aws/commit/f15217a139d543a5f8d1b6f13f74d6b433a2e218))

### [0.0.5](https://github.com/padok-team/yatas-aws/compare/v0.0.4...v0.0.5) (2022-09-09)


### Features

* **renovate:** added pretty changelog ([fbd8125](https://github.com/padok-team/yatas-aws/commit/fbd812557d77a89c98d6af2225fd54301a394896))

### [0.0.4](https://github.com/padok-team/yatas-aws/compare/v0.0.3...v0.0.4) (2022-09-09)


### Features

* **license:** added ([59ca94b](https://github.com/padok-team/yatas-aws/commit/59ca94b7c19b03c61849a890f4fa6586bc125306))
* **refacto:** changed package ([2f06d96](https://github.com/padok-team/yatas-aws/commit/2f06d9679f8bc1518b5c6105ab9db94c520c44cc))
* **renovate:** added config ([6215cf4](https://github.com/padok-team/yatas-aws/commit/6215cf4724a5b03bedc0ea3fd1b3ccc53dfb8600))
* **standardversion:** changed version ([9712c0b](https://github.com/padok-team/yatas-aws/commit/9712c0b78417e46602c087a9cc7b41c838e9eed7))
