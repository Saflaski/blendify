export function Privacy() {
  return (
    <div className="w-full flex justify-center py-12 px-4">
      <div className="w-full max-w-3xl ring-2 ring-black  px-8 py-10">
        <div className="text-black font-[Roboto_Mono] space-y-10">
          <header className="text-center space-y-2">
            <h1 className="text-3xl font-semibold">Privacy Policy</h1>
            <p className="text-sm">
              <strong>Last updated:</strong> January 1, 2026
            </p>
          </header>

          <section className="space-y-4">
            <h2 className="text-xl font-semibold border-b border-black pb-2">
              Overview
            </h2>
            <p>
              <strong>Blendify</strong> (“the Site”) is operated by{" "}
              <strong>Sheikh Sabeeh Islam</strong> and is a public website
              available at <strong>blendify.sabeehislam.com</strong>.
            </p>
            <p>
              We collect the minimum amount of data necessary to operate the
              Site.
            </p>
          </section>

          <section className="space-y-6">
            <h2 className="text-xl font-semibold border-b border-black pb-2">
              Information We Collect
            </h2>

            <div className="space-y-3">
              <h3 className="text-lg font-semibold">
                Account &amp; Login Data
              </h3>
              <p>
                Users authenticate via a third-party provider (
                <strong>Last.fm</strong>).
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li>
                  We store <strong>only a username</strong>
                </li>
                <li>
                  We do <strong>not</strong> store email addresses, passwords,
                  or credentials
                </li>
                <li>
                  A <strong>one-time authentication token</strong> is used
                  during first login and discarded immediately
                </li>
              </ul>
            </div>

            <div className="space-y-3">
              <h3 className="text-lg font-semibold">Cookies</h3>
              <p>We use a single essential cookie:</p>
              <ul className="list-disc pl-6 space-y-2">
                <li>
                  <strong>Session ID cookie</strong>
                  <ul className="list-disc pl-6 mt-2 space-y-1">
                    <li>Purpose: Maintain user sessions</li>
                    <li>Expiration: 7 days</li>
                  </ul>
                </li>
              </ul>
              <p>No analytics, advertising, or tracking cookies are used.</p>
            </div>

            <div className="space-y-3">
              <h3 className="text-lg font-semibold">
                Server &amp; Application Logs
              </h3>
              <p>
                For security and operational purposes, we log limited events,
                including login activity.
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li>Username</li>
                <li>Timestamp</li>
                <li>IP address</li>
                <li>User agent</li>
              </ul>
              <p>
                Logs are used only for security, debugging, and system
                reliability. They are <strong>not</strong> used for analytics,
                profiling, or marketing and are retained only as long as
                necessary.
              </p>
            </div>
          </section>

          <section className="space-y-4">
            <h2 className="text-xl font-semibold border-b border-black pb-2">
              How We Use Data
            </h2>
            <ul className="list-disc pl-6 space-y-2">
              <li>Authenticate users</li>
              <li>Maintain sessions</li>
              <li>Secure and operate the Site</li>
            </ul>
            <p>
              We do not sell or share personal data for commercial purposes.
            </p>
          </section>

          <section className="space-y-4">
            <h2 className="text-xl font-semibold border-b border-black pb-2">
              Third Parties
            </h2>
            <ul className="list-disc pl-6 space-y-2">
              <li>
                <strong>Authentication:</strong> Last.fm
              </li>
              <li>
                <strong>Hosting:</strong> Hetzner (VPS)
              </li>
            </ul>
            <p>
              Error logging is handled internally and not shared with
              third-party monitoring services.
            </p>
          </section>

          <section className="space-y-4">
            <h2 className="text-xl font-semibold border-b border-black pb-2">
              Your Rights (GDPR / CCPA)
            </h2>
            <ul className="list-disc pl-6 space-y-2">
              <li>Access your data</li>
              <li>Request deletion</li>
              <li>Request correction</li>
            </ul>
            <p>
              You can <strong>view or delete your data directly</strong> via the{" "}
              <strong>Account</strong> page when logged in.
            </p>
            <p>If you need additional assistance, you may contact us.</p>
          </section>

          <section className="space-y-2">
            <h2 className="text-xl font-semibold border-b border-black pb-2">
              Contact
            </h2>
            <p>For privacy-related questions:</p>
            <p>
              <strong>Email:</strong>{" "}
              <a className="underline" rel="noopener">
                blendify@sabeehislam.com
              </a>
            </p>
          </section>
        </div>
      </div>
    </div>
  );
}
