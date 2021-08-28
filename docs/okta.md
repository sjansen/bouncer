# Creating SP Configuration

1. Directory > Groups > Add Group
 - Name: bouncer-admin
 - Click Assign People and at least one user
1. Applications > Create App Integration
 - Do not use a template
 - Sign on method: SAML 2.0
1. General Settings
 - App name: bouncer
1. SAML Settings
 - Single sign on URL: https://example.com/b/saml/acs
 - Use this for Recipient URL and Destination URL: true
 - Allow this app to request other SSO URLs: false
 - Audience URI (SP Entity ID): bouncer
 - Default RelayState: (blank)
 - Name ID format: Unspecified
 - Application username: email
 - Update application username on: Create and update
1. Attribute Statements
 - firstName Unspecified user.firstName
 - lastName Unspecified user.lastName
 - email Unspecified user.email
 - username Unspecified String.substringBefore(user.email, "@") 
1. Group Attribute Statements
 - roles Unspecified Starts with: bouncer- 
1. Save configuration.
1. Copy "Identity Provider metadata" link
 - Add URL to SSM as `SAML_METADATA_URL`
1. Generate a PEM formatted RSA key pair
 - Add certificate to SSM as `SAML_CERTIFICATE`
 - Add private key to SSM as `SAML_PRIVATE_KEY`
1. Click on Assignments tab
 - Assign > Assign to Groups
 - Assign bouncer-admin
