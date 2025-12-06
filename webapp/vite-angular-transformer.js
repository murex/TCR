/*
Copyright (c) 2024 Murex

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

const fs = require('fs');
const path = require('path');

/**
 * Simple Angular component transformer for Vitest
 * This plugin transforms Angular components by:
 * 1. Reading templateUrl and styleUrl files
 * 2. Replacing them with inline template and styles
 * 3. Making components work in Vitest environment
 */
function angularTransformer() {
  return {
    name: 'angular-transformer',
    transform(code, id) {
      // Only transform .ts files that are not test files
      if (!id.endsWith('.ts') || id.includes('.spec.ts') || id.includes('.test.ts')) {
        return null;
      }

      // Look for Angular components with templateUrl or styleUrl
      const templateUrlRegex = /templateUrl:\s*['"`]([^'"`]+)['"`]/g;
      const styleUrlRegex = /styleUrl:\s*['"`]([^'"`]+)['"`]/g;
      const styleUrlsRegex = /styleUrls:\s*\[([^\]]+)\]/g;

      let transformedCode = code;
      let hasTransformations = false;

      // Transform templateUrl to template
      transformedCode = transformedCode.replace(templateUrlRegex, (match, templatePath) => {
        try {
          const fullTemplatePath = path.resolve(path.dirname(id), templatePath);
          if (fs.existsSync(fullTemplatePath)) {
            const templateContent = fs.readFileSync(fullTemplatePath, 'utf-8');
            // Escape backticks and ${} in template content
            const escapedContent = templateContent
              .replace(/\\/g, '\\\\')
              .replace(/`/g, '\\`')
              .replace(/\$\{/g, '\\${');
            hasTransformations = true;
            return `template: \`${escapedContent}\``;
          }
        } catch (error) {
          console.warn(`Failed to read template file: ${templatePath}`, error.message);
        }
        return match;
      });

      // Transform styleUrl to styles (single style)
      transformedCode = transformedCode.replace(styleUrlRegex, (match, stylePath) => {
        try {
          const fullStylePath = path.resolve(path.dirname(id), stylePath);
          if (fs.existsSync(fullStylePath)) {
            const styleContent = fs.readFileSync(fullStylePath, 'utf-8');
            // Escape backticks and ${} in style content
            const escapedContent = styleContent
              .replace(/\\/g, '\\\\')
              .replace(/`/g, '\\`')
              .replace(/\$\{/g, '\\${');
            hasTransformations = true;
            return `styles: [\`${escapedContent}\`]`;
          }
        } catch (error) {
          console.warn(`Failed to read style file: ${stylePath}`, error.message);
        }
        return match;
      });

      // Transform styleUrls to styles (array of styles)
      transformedCode = transformedCode.replace(styleUrlsRegex, (match, styleUrlsContent) => {
        try {
          // Extract individual style URLs from the array
          const styleUrls = styleUrlsContent
            .split(',')
            .map(url => url.trim().replace(/['"`]/g, ''))
            .filter(url => url);

          const styles = [];
          for (const styleUrl of styleUrls) {
            const fullStylePath = path.resolve(path.dirname(id), styleUrl);
            if (fs.existsSync(fullStylePath)) {
              const styleContent = fs.readFileSync(fullStylePath, 'utf-8');
              // Escape backticks and ${} in style content
              const escapedContent = styleContent
                .replace(/\\/g, '\\\\')
                .replace(/`/g, '\\`')
                .replace(/\$\{/g, '\\${');
              styles.push(`\`${escapedContent}\``);
            }
          }

          if (styles.length > 0) {
            hasTransformations = true;
            return `styles: [${styles.join(', ')}]`;
          }
        } catch (error) {
          console.warn(`Failed to read style files: ${styleUrlsContent}`, error.message);
        }
        return match;
      });

      return hasTransformations ? { code: transformedCode, map: null } : null;
    }
  };
}

module.exports = { angularTransformer };
