// Polyfill for Google Closure Library (minimal implementation needed for protobuf)
import * as jspb from 'google-protobuf';

// Define global if it doesn't exist (browser environment)
if (typeof window !== 'undefined') {
  // @ts-ignore
  if (typeof window.global === 'undefined') {
    // @ts-ignore
    window.global = window;
  }
  // @ts-ignore
  if (typeof global === 'undefined') {
    // @ts-ignore
    globalThis.global = window;
  }
  
  // Make jspb available globally for generated protobuf code
  // @ts-ignore
  window.jspb = jspb;
}

// Initialize goog for Google Closure Library
// @ts-ignore
if (typeof window !== 'undefined' && typeof window.goog === 'undefined') {
  // @ts-ignore
  window.goog = {};
  // @ts-ignore
  window.goog.exportSymbol = function(name: string, value: any, object: any) {
    if (!object) object = window;
    const parts = name.split('.');
    let cur = object;
    for (let i = 0; i < parts.length - 1; i++) {
      const part = parts[i];
      if (!cur[part]) cur[part] = {};
      cur = cur[part];
    }
    cur[parts[parts.length - 1]] = value;
  };
  // @ts-ignore
  window.goog.inherits = function(childCtor: any, parentCtor: any) {
    // Proper prototype chain setup
    childCtor.superClass_ = parentCtor.prototype;
    childCtor.prototype = Object.create(parentCtor.prototype);
    childCtor.prototype.constructor = childCtor;
  };
  // @ts-ignore
  window.goog.DEBUG = false;
}

export {};

