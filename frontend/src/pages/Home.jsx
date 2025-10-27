import { Link } from 'react-router-dom';

function Home() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-500 to-purple-600">
      <div className="bg-white rounded-lg shadow-2xl p-8 max-w-md w-full text-center">
        <h1 className="text-4xl font-bold text-gray-800 mb-4">Welcome to GoChatApp</h1>
        <p className="text-gray-600 mb-8">Connect and chat with others in real-time</p>
        <div className="space-y-4">
          <Link
            to="/login"
            className="block w-full bg-blue-600 hover:bg-blue-700 text-white font-semibold py-3 px-6 rounded-lg transition duration-200"
          >
            Login
          </Link>
          <Link
            to="/register"
            className="block w-full bg-purple-600 hover:bg-purple-700 text-white font-semibold py-3 px-6 rounded-lg transition duration-200"
          >
            Register
          </Link>
        </div>
      </div>
    </div>
  );
}

export default Home;
