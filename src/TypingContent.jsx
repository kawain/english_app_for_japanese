import VolumeControl from './components/VolumeControl.jsx'

function TypingContent ({ volume, onVolumeChange }) {
  return (
    <>
      <h1>タイピング</h1>
      <VolumeControl volume={volume} onVolumeChange={onVolumeChange} />
    </>
  )
}
export default TypingContent
