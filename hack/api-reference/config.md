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


