import {
  type PublicKeyCredentialCreationOptionsJSON,
  type RegistrationResponseJSON,
  startRegistration,
  WebAuthnError,
} from "@simplewebauthn/browser";

export type PublicKey = PublicKeyCredentialCreationOptionsJSON;
export type RegistrationResponse = RegistrationResponseJSON;

/**
 * Start the registration process for a new passkey.
 *
 * @param publicKey Output from **@simplewebauthn/server**'s `generateRegistrationOptions()`
 * @throws Error if the registration fails
 */
export const createPasskey = async (
  publicKey: PublicKey,
): Promise<RegistrationResponse> => {
  try {
    const res = await startRegistration({ optionsJSON: publicKey });
    return res;
  } catch (error) {
    if (error instanceof WebAuthnError) {
      console.error("WebAuthnError: Failed to create passkey", error);
      throw error;
    } else {
      throw error;
    }
  }
};
