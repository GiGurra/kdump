use aes_gcm::{Aes256Gcm, Nonce};
use aes_gcm::aes::cipher::consts::U12;
use aes_gcm::aead::{Aead, NewAead};

use rand_chacha;
use rand::prelude::*;

pub fn encrypt(input: &str, key: &str) -> Vec<u8> {
    // https://docs.rs/aes-gcm/latest/aes_gcm/
    use rand_chacha::ChaCha20Rng;

    let mut rng = ChaCha20Rng::from_entropy();

    let nonce_bytes_as_u32: Vec<u32> = (0..3).map(|_| rng.next_u32()).collect();
    let nonce_bytes: Vec<u8> = nonce_bytes_as_u32.iter().flat_map(|x| x.to_be_bytes()).collect();

    let nonce = Nonce::<U12>::from_slice(&nonce_bytes);
    let cipher = Aes256Gcm::new_from_slice(key.as_bytes()).unwrap();
    let ciphertext = cipher.encrypt(nonce, input.as_bytes()).unwrap();

    return ciphertext;
}
