<p>Packages:</p>
<ul>
<li>
<a href="#config.ubuntu.os.extensions.gardener.cloud%2fv1alpha1">config.ubuntu.os.extensions.gardener.cloud/v1alpha1</a>
</li>
</ul>

<h2 id="config.ubuntu.os.extensions.gardener.cloud/v1alpha1">config.ubuntu.os.extensions.gardener.cloud/v1alpha1</h2>
<p>

</p>

<h3 id="aptarchive">APTArchive
</h3>


<p>
(<em>Appears on:</em><a href="#aptconfig">APTConfig</a>)
</p>

<p>

</p>

<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>

<tr>
<td>
<code>arches</code></br>
<em>
<a href="#architecture">Architecture</a> array
</em>
</td>
<td>
<p></p>
</td>
</tr>
<tr>
<td>
<code>uri</code></br>
<em>
string
</em>
</td>
<td>
<p></p>
</td>
</tr>
<tr>
<td>
<code>search</code></br>
<em>
string array
</em>
</td>
<td>
<p></p>
</td>
</tr>
<tr>
<td>
<code>searchDNS</code></br>
<em>
boolean
</em>
</td>
<td>
<p></p>
</td>
</tr>

</tbody>
</table>


<h3 id="aptconfig">APTConfig
</h3>


<p>
(<em>Appears on:</em><a href="#extensionconfig">ExtensionConfig</a>)
</p>

<p>

</p>

<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>

<tr>
<td>
<code>preserveSourcesList</code></br>
<em>
boolean
</em>
</td>
<td>
<p></p>
</td>
</tr>
<tr>
<td>
<code>primary</code></br>
<em>
<a href="#aptarchive">APTArchive</a> array
</em>
</td>
<td>
<p></p>
</td>
</tr>
<tr>
<td>
<code>security</code></br>
<em>
<a href="#aptarchive">APTArchive</a> array
</em>
</td>
<td>
<p></p>
</td>
</tr>

</tbody>
</table>


<h3 id="aptrepository">AptRepository
</h3>


<p>
(<em>Appears on:</em><a href="#extensionconfig">ExtensionConfig</a>)
</p>

<p>
AptRepository describes an additional apt repository to configure via
cloud-init.
</p>

<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>

<tr>
<td>
<code>name</code></br>
<em>
string
</em>
</td>
<td>
<p>Name is a unique name for the apt source.</p>
</td>
</tr>
<tr>
<td>
<code>uri</code></br>
<em>
string
</em>
</td>
<td>
<p>URI is the base URI of the apt repository. For the docker repository this<br />is typically https://download.docker.com/linux/ubuntu.</p>
</td>
</tr>
<tr>
<td>
<code>key</code></br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Key is the ASCII-armored GPG key used to sign the repository. If empty,<br />the repository is configured without GPG signature verification.</p>
</td>
</tr>
<tr>
<td>
<code>keyUrl</code></br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>KeyURL is the URL to download the GPG key from. The key is downloaded<br />via cloud-init write_files to /etc/apt/keyrings/<name>.<keyFormat> and<br />referenced with signed-by. Mutually exclusive with Key.</p>
</td>
</tr>
<tr>
<td>
<code>keyFormat</code></br>
<em>
<a href="#keyformat">KeyFormat</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>KeyFormat specifies the format of the GPG key served by KeyURL. It<br />determines the file extension of the key written by cloud-init, which<br />must match the key's content: "asc" for ASCII-armored keys or "gpg" for<br />binary keyrings. Defaults to "asc".</p>
</td>
</tr>
<tr>
<td>
<code>suite</code></br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Suite is the apt suite to use. Defaults to "$RELEASE" which cloud-init<br />substitutes with the release codename.</p>
</td>
</tr>
<tr>
<td>
<code>components</code></br>
<em>
string array
</em>
</td>
<td>
<em>(Optional)</em>
<p>Components is the list of apt components to use. Defaults to ["stable"].</p>
</td>
</tr>

</tbody>
</table>


<h3 id="architecture">Architecture
</h3>
<p><em>Underlying type: string</em></p>


<p>
(<em>Appears on:</em><a href="#aptarchive">APTArchive</a>)
</p>

<p>

</p>


<h3 id="daemon">Daemon
</h3>
<p><em>Underlying type: string</em></p>


<p>
(<em>Appears on:</em><a href="#ntpconfig">NTPConfig</a>)
</p>

<p>

</p>


<h3 id="dependencyconfig">DependencyConfig
</h3>


<p>
(<em>Appears on:</em><a href="#extensionconfig">ExtensionConfig</a>)
</p>

<p>
DependencyConfig describes an apt package to install.
</p>

<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>

<tr>
<td>
<code>name</code></br>
<em>
string
</em>
</td>
<td>
<p>Name is the name of the apt package.</p>
</td>
</tr>
<tr>
<td>
<code>version</code></br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Version is the exact apt package version to install. If empty, the latest<br />available version is installed.</p>
</td>
</tr>
<tr>
<td>
<code>ubuntuVersion</code></br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>UbuntuVersion optionally restricts this dependency to a specific Ubuntu<br />version (e.g. "22.04"). The value is matched against VERSION_ID from<br />/etc/os-release.</p>
</td>
</tr>
<tr>
<td>
<code>ubuntuBuildSerial</code></br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>UbuntuBuildSerial optionally restricts this dependency to a specific<br />Ubuntu build serial (e.g. "20250725"). The value is matched against the<br />serial field in /etc/cloud/build.info.</p>
</td>
</tr>
<tr>
<td>
<code>hold</code></br>
<em>
boolean
</em>
</td>
<td>
<em>(Optional)</em>
<p>Hold, if true, marks the package on apt hold after installation so that<br />apt upgrade does not update it.</p>
</td>
</tr>

</tbody>
</table>


<h3 id="extensionconfig">ExtensionConfig
</h3>


<p>
ExtensionConfig is the configuration for the os-ubuntu extension.
</p>

<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>

<tr>
<td>
<code>ntp</code></br>
<em>
<a href="#ntpconfig">NTPConfig</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>NTP to configure either systemd-timesyncd or ntpd</p>
</td>
</tr>
<tr>
<td>
<code>disableUnattendedUpgrades</code></br>
<em>
boolean
</em>
</td>
<td>
<em>(Optional)</em>
<p>DisableUnattendedUpgrades to disable unattended upgrades in ubuntu</p>
</td>
</tr>
<tr>
<td>
<code>aptRepositories</code></br>
<em>
<a href="#aptrepository">AptRepository</a> array
</em>
</td>
<td>
<em>(Optional)</em>
<p>AptRepositories is the list of additional apt repositories to configure<br />via cloud-init.</p>
</td>
</tr>
<tr>
<td>
<code>dependencies</code></br>
<em>
<a href="#dependencyconfig">DependencyConfig</a> array
</em>
</td>
<td>
<em>(Optional)</em>
<p>Dependencies is the list of apt packages to install on the node. If empty,<br />a default set of unpinned packages is installed. Each dependency may<br />optionally target a specific Ubuntu version and/or build serial.</p>
</td>
</tr>
<tr>
<td>
<code>apt</code></br>
<em>
<a href="#aptconfig">APTConfig</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Mirror to set custom Ubuntu mirror</p>
</td>
</tr>

</tbody>
</table>


<h3 id="keyformat">KeyFormat
</h3>
<p><em>Underlying type: string</em></p>


<p>
(<em>Appears on:</em><a href="#aptrepository">AptRepository</a>)
</p>

<p>
KeyFormat specifies the format of a GPG key. It determines the file
extension of the key file written by cloud-init, which must match the key's
content.
</p>


<h3 id="ntpconfig">NTPConfig
</h3>


<p>
(<em>Appears on:</em><a href="#extensionconfig">ExtensionConfig</a>)
</p>

<p>
NTPConfig General NTP Config for either systemd-timesyncd or ntpd
</p>

<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>

<tr>
<td>
<code>daemon</code></br>
<em>
<a href="#daemon">Daemon</a>
</em>
</td>
<td>
<p>Daemon One of either systemd-timesyncd or ntp</p>
</td>
</tr>
<tr>
<td>
<code>ntpd</code></br>
<em>
<a href="#ntpdconfig">NTPDConfig</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>NTPD to configure the ntpd client</p>
</td>
</tr>

</tbody>
</table>


<h3 id="ntpdconfig">NTPDConfig
</h3>


<p>
(<em>Appears on:</em><a href="#ntpconfig">NTPConfig</a>)
</p>

<p>
NTPDConfig is the struct used in the ntp-config.conf.tpl template file
</p>

<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>

<tr>
<td>
<code>servers</code></br>
<em>
string array
</em>
</td>
<td>
<p>Servers List of ntp servers</p>
</td>
</tr>
<tr>
<td>
<code>interfaces</code></br>
<em>
string array
</em>
</td>
<td>
<p>Interfaces for ntpd to bind to. Can be more than one.</p>
</td>
</tr>

</tbody>
</table>


