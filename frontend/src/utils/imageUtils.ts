/**
 * Image utility functions for handling base64 images
 */

/**
 * Converts a raw base64 string to a data URI that can be used in img src
 * @param base64String - Raw base64 string without data URI prefix
 * @param mimeType - MIME type of the image (default: image/png)
 * @returns Data URI string ready for use in img src attribute
 */
export function toDataUri(base64String: string | null | undefined, mimeType: string = 'image/png'): string {
  if (!base64String) {
    return '';
  }

  // If already has data URI prefix, return as is
  if (base64String.startsWith('data:')) {
    return base64String;
  }

  // If it's a URL (http/https), return as is
  if (base64String.startsWith('http://') || base64String.startsWith('https://')) {
    return base64String;
  }

  // Add data URI prefix to raw base64 string
  return `data:${mimeType};base64,${base64String}`;
}

/**
 * Checks if a string is a valid base64 string
 * @param str - String to check
 * @returns true if the string appears to be base64 encoded
 */
export function isBase64(str: string): boolean {
  if (!str || str.length === 0) {
    return false;
  }

  // Check if already has data URI prefix
  if (str.startsWith('data:')) {
    return true;
  }

  // Check if it's a URL
  if (str.startsWith('http://') || str.startsWith('https://')) {
    return false;
  }

  // Check if it matches base64 pattern
  const base64Regex = /^[A-Za-z0-9+/]*={0,2}$/;
  return base64Regex.test(str);
}
