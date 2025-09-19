import { SeriesList } from '@/components/SeriesList';

export default function Home(): React.JSX.Element {
  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white shadow-sm border-b">
        <div
          className="
          w-full max-w-sm mx-auto px-4 py-4
          sm:max-w-md sm:px-6 sm:py-5
          md:max-w-lg md:px-8 md:py-6
        "
        >
          <div
            className="flex flex-col items-center space-y-4
            sm:flex-row sm:justify-between sm:space-y-0
          "
          >
            <h1
              className="
              text-xl font-bold text-gray-900 text-center
              sm:text-2xl md:text-3xl
            "
            >
              Spark Park Cricket
            </h1>
          </div>
        </div>
      </header>

      <main
        className="
        w-full max-w-sm mx-auto px-4 py-6
        sm:max-w-md sm:px-6 sm:py-8
        md:max-w-lg md:px-8 md:py-10
        lg:max-w-xl lg:py-12
      "
      >
        <div className="mb-8 text-center">
          <h2
            className="
            text-xl font-bold text-gray-900 mb-4
            sm:text-2xl md:text-3xl
          "
          >
            Welcome to Spark Park Cricket
          </h2>
          <p
            className="
            text-sm text-gray-600
            sm:text-base md:text-lg
          "
          >
            Manage your cricket tournaments, matches, and teams with our
            comprehensive tournament management system.
          </p>
        </div>

        <SeriesList />
      </main>

      <footer className="bg-white border-t">
        <div
          className="
          w-full max-w-sm mx-auto py-4 px-4
          sm:max-w-md sm:px-6 sm:py-6
          md:max-w-lg md:px-8
        "
        >
          <p
            className="text-center text-gray-500 text-xs
            sm:text-sm
          "
          >
            © 2024 Spark Park Cricket. All rights reserved.
          </p>
        </div>
      </footer>
    </div>
  );
}
