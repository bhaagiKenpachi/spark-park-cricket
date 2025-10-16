# Cloudflare CSP Configuration for Google AdSense

> ✅ **Status**: Issue resolved using Solution 2 (Disable Cloudflare's CSP)

## Issue
Cloudflare is applying its own Content Security Policy (CSP) that blocks Google AdSense scripts, overriding the Next.js CSP configuration.

## Error Messages
```
Refused to load the script 'https://pagead2.googlesyndication.com/pagead/js/adsbygoogle.js' 
because it violates the following Content Security Policy directive: 
"script-src 'self' https://static.cloudflareinsights.com https://challenges.cloudflare.com 'unsafe-inline' 'unsafe-eval'"
```

## Solution: Configure Cloudflare to Allow AdSense

### Option 1: Disable Cloudflare's CSP (Recommended) ✅ WORKING

1. **Log in to Cloudflare Dashboard**
   - Go to https://dash.cloudflare.com
   - Select your domain (e.g., `dojima.foundation`)

2. **Navigate to Security Settings**
   - Click on **"Security"** in the left sidebar
   - Go to **"Settings"** tab

3. **Find Content Security Policy**
   - Look for **"Browser Integrity Check"** or **"Security Headers"**
   - If there's a CSP setting, disable it or modify it

4. **Alternative: Use Transform Rules**
   - Go to **"Rules"** → **"Transform Rules"** → **"Modify Response Header"**
   - Create a new rule to remove or modify the CSP header
   - Rule name: "Remove CSP for AdSense"
   - When incoming requests match: `(http.host eq "spark-park-dev.dojima.foundation")`
   - Then: **Remove** header named `Content-Security-Policy`
   - Deploy the rule

### Option 2: Modify Cloudflare's CSP to Include AdSense Domains

If you want to keep Cloudflare's CSP but allow AdSense:

1. **Go to Transform Rules**
   - Navigate to **"Rules"** → **"Transform Rules"** → **"Modify Response Header"**

2. **Create CSP Override Rule**
   - Rule name: "CSP for AdSense"
   - When incoming requests match: `(http.host eq "spark-park-dev.dojima.foundation")`
   - Then: **Set static** header named `Content-Security-Policy`
   - Value:
   ```
   default-src 'self' 'unsafe-inline' 'unsafe-eval' data: blob:; script-src 'self' 'unsafe-inline' 'unsafe-eval' https://pagead2.googlesyndication.com https://*.google.com https://*.adtrafficquality.google https://www.googletagmanager.com https://static.cloudflareinsights.com https://challenges.cloudflare.com; connect-src 'self' https://spark-park-dev.dojima.foundation https://ochhmsslirapqqzcgvek.supabase.co https://api.whatsapp.com https://api.iconify.design https://pagead2.googlesyndication.com https://*.google.com https://*.google-analytics.com https://*.adtrafficquality.google https://cloudflareinsights.com; img-src 'self' data: https: blob:; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; font-src 'self' data: https://fonts.gstatic.com; frame-src 'self' https://googleads.g.doubleclick.net https://tpc.googlesyndication.com https://*.adtrafficquality.google https://www.google.com;
   ```

3. **Deploy and Test**
   - Save the rule
   - Wait 1-2 minutes for propagation
   - Clear browser cache and test

### Option 3: Page Rules (Legacy)

If Transform Rules aren't available:

1. **Go to Page Rules**
   - Navigate to **"Rules"** → **"Page Rules"**

2. **Create New Page Rule**
   - URL: `spark-park-dev.dojima.foundation/*`
   - Add Setting: **"Disable Security"** or **"Browser Integrity Check: Off"**
   - Save and Deploy

## Verification

After configuring Cloudflare:

1. **Clear Browser Cache**
   - Chrome: DevTools → Network → "Disable cache" checkbox
   - Or use Incognito/Private browsing

2. **Check Response Headers**
   - Open DevTools (F12)
   - Go to Network tab
   - Refresh the page
   - Click on the document request
   - Check the "Response Headers" section
   - Verify that the CSP header includes AdSense domains

3. **Verify AdSense Loads**
   - Check Console for CSP errors (should be gone)
   - Verify AdSense script loads successfully
   - Check that ad slots are rendering

## Testing Domains

Apply these configurations to:
- **Dev**: `spark-park-dev.dojima.foundation`
- **Prod**: `spark-park.dojima.foundation`

## Important Notes

1. **CSP Priority**: Cloudflare's CSP headers override Next.js headers when both are present
2. **Propagation Time**: Allow 1-2 minutes for Cloudflare changes to propagate globally
3. **Cache**: Clear browser cache after making changes
4. **Multiple Domains**: Create separate rules for dev and production domains

## Troubleshooting

### Issue: CSP errors still appear after configuration
**Solution**: 
- Wait 2-3 minutes for Cloudflare to propagate changes
- Clear browser cache completely
- Try incognito/private browsing mode
- Check that the rule is enabled in Cloudflare dashboard

### Issue: AdSense still not loading
**Solution**:
- Verify the publisher ID is correct: `ca-pub-5474524579770573`
- Check that ad slots are created in AdSense dashboard
- Ensure the domain is added to AdSense account
- Wait for AdSense approval (can take 24-48 hours)

### Issue: Cloudflare blocking other scripts
**Solution**:
- Add those domains to the CSP whitelist
- Use wildcards carefully (e.g., `https://*.google.com`)

## Additional Resources

- [Cloudflare Transform Rules Documentation](https://developers.cloudflare.com/rules/transform/)
- [Google AdSense CSP Requirements](https://support.google.com/adsense/answer/12171612)
- [Content Security Policy (CSP) Guide](https://developer.mozilla.org/en-US/docs/Web/HTTP/CSP)

